package pkgep

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"strings"

// 	"gorm.io/gorm"
// )

// // handle gorm errors
// func GormErr(err error) error {
// 	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
// 		fmt.Printf("err == gorm.ErrCheckConstraintViolated %s\n", err.Error())
// 		return schema.SchemaRepoError{
// 			Code:    schema.DbConflictErr,
// 			Message: "Invalid option ID",
// 		}
// 	}

// 	if errors.Is(err, gorm.ErrForeignKeyViolated) {
// 		fmt.Printf("err == gorm.ErrForeignKeyViolated %s\n", err.Error())
// 		return schema.SchemaRepoError{
// 			Code:    schema.DbInternalErr,
// 			Message: "Invalid option ID",
// 		}
// 	}

// 	if strings.Contains(err.Error(), "disaster_name_th_idx") || strings.Contains(err.Error(), "disaster_name_en_idx") {

// 		re := regexp.MustCompile(`\((\d+(?:,\s*\d+)*)\)`)
// 		matches := re.FindAllStringSubmatch(err.Error(), -1)

// 		conflictID := ""

// 		if len(matches) > 0 && len(matches[0]) > 1 {
// 			conflictID = matches[0][1]
// 		}

// 		return schema.SchemaRepoError{
// 			Code:    schema.DbConflictErr,
// 			Message: fmt.Sprintf("Conflict, Question with ID %s have already been answered", strings.Split(conflictID, ",")[0]),
// 		}
// 	}
// }
