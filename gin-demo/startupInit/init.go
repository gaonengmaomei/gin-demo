package startupInit

import (
	"fmt"
	"gin-demo/constant"
	"gin-demo/databases"
	"gin-demo/models"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

func Load() {
	loadModels()
}

func loadModels() {
	modelsFilePath, _ := filepath.Abs(fmt.Sprintf("./%s", constant.ModelsFileName))

	dir, err := os.Open(modelsFilePath)
	if err != nil {
		panic("models文件夹不存在")
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		panic("models不存在文件")
	}

	modelMap := make(map[string]*models.ModelStruct, 100)

	for _, file := range files {
		if !checkModel(file.Name()) {
			continue
		}

		fullPath := fmt.Sprintf("%s/%s", modelsFilePath, file.Name())
		fset := token.NewFileSet()
		astfile, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
		if err != nil {
			continue
		}
		if astfile.Name.Name != constant.ModelsFileName {
			continue
		}
		//ast.Print(fset, astfile)

		//添加字段
		for _, v := range astfile.Decls {
			if decl, ok := v.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					if tp, ok := spec.(*ast.TypeSpec); ok {
						if stp, ok := tp.Type.(*ast.StructType); ok {
							if !stp.Struct.IsValid() {
								continue
							}
							modelStruct := models.ModelStruct{}

							for _, field := range stp.Fields.List {
								elementTypeReflect := getElementType(field)
								if elementTypeReflect == nil {
									panic(fmt.Sprintf("不支持的数据类型,%v", field))
								}

								if field.Tag != nil {
									tag := strings.ReplaceAll(field.Tag.Value, "`", "")
									// 新增成员
									modelStruct.StructField = append(modelStruct.StructField, reflect.StructField{
										Name: field.Names[0].Name,
										Type: elementTypeReflect,
										Tag:  reflect.StructTag(tag), // 将string类型强制转换为reflect.StructTag类型
									})
								} else {
									// 新增成员
									modelStruct.StructField = append(modelStruct.StructField, reflect.StructField{
										Name: field.Names[0].Name,
										Type: elementTypeReflect,
									})
								}

							}
							modelMap[tp.Name.Name] = &modelStruct
						}
					}
				}
			}

			if fun, ok := v.(*ast.FuncDecl); ok && fun.Name.Name == "TableName" {
				if fun.Recv != nil && len(fun.Recv.List) > 0 {
					if Ident, ok := fun.Recv.List[0].Type.(*ast.Ident); ok {
						modelName := Ident.Name
						if oneModelStruct, ok := modelMap[modelName]; ok {
							tableName := fun.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value
							tableName = strings.ReplaceAll(tableName, "\"", "")
							oneModelStruct.TableName = tableName
						}
					}
				}
			}
		}
	}

	databases.AutoMigrate(modelMap)
}

// 这里来解析models成员的类型，可以根据实际情况扩展
func getElementType(field *ast.Field) reflect.Type {
	// 判断成员类型，可根据实际情况扩展
	var (
		elementType        string
		elementTypeReflect reflect.Type
	)

	if _, ok := field.Type.(*ast.Ident); ok {
		elementType = field.Type.(*ast.Ident).Name
	}
	// 如果是time.Time会被ast识别为ast.SelectorExpr
	if _, ok := field.Type.(*ast.SelectorExpr); ok {
		elementType = field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + field.Type.(*ast.SelectorExpr).Sel.Name
	}

	switch elementType {
	case "int":
		elementTypeReflect = reflect.TypeOf(1)
	case "int8":
		elementTypeReflect = reflect.TypeOf(int8(1))
	case "int16":
		elementTypeReflect = reflect.TypeOf(int16(1))
	case "int32":
		elementTypeReflect = reflect.TypeOf(int32(1))
	case "int64":
		elementTypeReflect = reflect.TypeOf(int64(1))
	case "float32":
		elementTypeReflect = reflect.TypeOf(float32(1))
	case "float64":
		elementTypeReflect = reflect.TypeOf(float64(1))
	case "bool":
		elementTypeReflect = reflect.TypeOf(true)
	case "string":
		elementTypeReflect = reflect.TypeOf("")
	case "time.Time":
		elementTypeReflect = reflect.TypeOf(time.Now())
	}

	return elementTypeReflect
}

func checkModel(name string) bool {
	name = strings.Replace(name, ".go", "", -1)
	suffixLen := utf8.RuneCountInString(constant.ModelsFileNameSuffix)
	nameLen := utf8.RuneCountInString(name)
	return strings.Index(name, constant.ModelsFileNameSuffix) == (nameLen - suffixLen)
}
