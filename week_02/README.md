# Week 02 错误处理

## 作业题目

- 假设我们在数据库操作的时候，比如 DAO 层中遇到一个 `sql.ErrNoRows` 的时候，是否应该封装这个错误抛给上层？ 
- 为什么，应该怎么做？请写出代码。

## 作业分析

- 我们知道 `sql.ErrNoRows` 是标准库 `sql` 包中的一个预定义错误，既 Sentinel Error.
- 题目中举例提到是在 DAO 层遇到这个错误。下面根据不同的代码层级来分析，以最终得出结论。

## 1 - 设计一个类似于 `sqlx` 的基础包（基础层级）

- 假设我们开发一个类似于 [sqlx](https://github.com/jmoiron/sqlx) 的，功能较为基础的包。
- 它提供的是基于标准库 `sql` 包的拓展，更方便、更高级的一些 API 接口，例如 `StructScan` 方法。
- 但是它与 `sql` 一样，依然属于基础包的范畴，不对调用者的业务需求做过多的假设。
- 那么 `sqlx` 本来就应该与 `sql` 高度耦合，甚至直接兼容 `sql` 的许多 API.
- 在这种情况下，我认为 `sqlx` 的错误处理用以下两种方式都可以：
  - 直接保留 `sql` 的预定义错误并返回。
  - 使用 Go 1.13 的 `fmt.Errorf("... %w", err)` 做简单封装。

```go
package sqlx

import "database/sql"

// 对 sql.DB 进行拓展
type DB struct {
	*sql.DB
	someMoreInfo string
}

// 对 sql.Rows 进行拓展
type Rows struct {
	*sql.Rows
	someMoreInfo string
}

// 保留 sql 包的预定义错误
//
// 可以看到下面的 Queryx 方法中 sqlx.DB 直接调用了其内部 sql.DB 的 Query 方法
// 如果出现 sql.ErrNoRows 或其他预定义错误，则不进行任何特殊处理直接返回
// 虽然拓展了 sql 包中的 DB 和 Rows 类型，但因为依然和 sql 包高度耦合，所以可以沿用它的预定义错误
func (db *DB) Queryx(query string, args ...interface{}) (*Rows, error) {
	rs, err := db.DB.Query(query, args...)
	return &Rows{rs, db.someMoreInfo}, err     // 直接返回，调用者不能对结果有任何假设，必须先处理错误
}
```

```go
package sqlx

import (
	"database/sql"
	"fmt"
)

// 使用 Go 1.13 的 `fmt.Errorf` 做简单的封装
//
// 这样的话，基本和 sql 包保持兼容，但是可以包含一些额外信息
// 当调用者需要处理错误的时候，有多方方式：
//  - 直接做等值比较 if err == sqlx.ErrNoRows
//  - 用 errors.Is 或 errors.As 可以检查出 sqlx.ErrNoRows 封装错误
//  - 用 errors.Is 或 errors.As 也可以检查出 sql.ErrNoRows 原始错误
var ErrNoRows = fmt.Errorf("sqlx: %w", sql.ErrNoRows)
```

- 不建议在 `sqlx` 中使用 `pkg/errors`。因为其功能较为基础，应该交给调用者在业务代码中做错误封装。

## 2 - 在业务中使用（业务层级）

- 如果是在业务中遇到 `sql.ErrNoRows`，由于它属于标准包非常基础的错误:
  - 在业务中推荐使用 `pkg/errors` 中的 `errors.Wrap`、`errors.Wrapf` 方法进行封装并保存堆栈信息。
  - 之后的上游调用者应该直接返回，或使用 `errors.WithMessage`、`errors.WithMessagef` 进行不保存堆栈信息的再次封装。

```go
package business

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

type User struct {
	ID      int
	Email   string
	Name    string
}

// 这个接口最贴近 DAO 层。
type UserRepo interface {
	GetUserByID(id int) (*User, error)
}

// 业务层逻辑。
type SomeBusinessService struct {
	userRepo UserRepo
}

// 业务层逻辑。
func (srv *SomeBusinessService) ReadUserAndDoSomething(id int) error {
	user, err := srv.userRepo.GetUserByID(id)
	if errors.Is(err, sql.ErrNoRows) {                  // 用 errors.Is 判断是否为 sql.ErrNoRows 
		return errors.Wrap(err, "user not found")   // 业务层用 Wrap 封装 - 我们知道是找不到用户
	}
	if err != nil {
		return errors.Wrap(err, "generic error")    // 业务层用 Wrap 封装 - 此时是其他的错误
	} 

	// 假设后面还有大量的业务逻辑。 
	//... 
	log.Printf("%+v\n", user)
	return nil
}
```

## 3 - 结论

- 那么 DAO 到底更贴近哪个层级？具体情况具体分析。

### 贴近业务

- 假设这个 DAO 层功能比较复杂，含有一定的业务逻辑，更贴近业务层级，那么当出现 `sql.ErrNoRows` 的时候可以封装，甚至不处理。
  - 例如这个 DAO 是一个 ORM 形式的，它可以提前加载（Eager Load）一个用户的信用卡列表。
  - 那么当这个用户还没有添加信用卡的时候，数据库层某一时刻会出现 `sql.ErrNoRows`，但 ORM 不一定要报错。

```go
package user

import "github.com/go-gorm/gorm"

type CreditCard struct{}

type User struct {
	gorm.Model
	Email       string       `gorm:"column:email"`
	Name        string       `gorm:"column:name"`
	CreditCards []CreditCard `gorm:"foreignKey:user_id"` // sql.ErrNoRows -> 返回空列表即可
}
```

### 贴近基础

- 个人认为，更多时候 DAO 层功能应该还是比较简单，是对 `sql` 层的一个拓展 + 对 `model` 的引入。
- 它应该最贴近上面第二部分例子中的 `UserRepo` 接口，介于业务层级和基础层级之间，**但更贴近基础层级**。
- **根据以上分析，当 DAO 层里遇到 `sql.ErrNoRows` 的时候，以下两种方式都可以：**
  - **直接将它返回。**
  - **或使用 Go 1.13 的 `fmt.Errorf("... %w", err)` 做简单封装。**
