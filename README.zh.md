# goda

[English](README.md) | 中文

[![CI](https://github.com/iseki0/goda/workflows/CI/badge.svg)](https://github.com/iseki0/goda/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/iseki0/goda.svg)](https://pkg.go.dev/github.com/iseki0/goda)
[![Go Report Card](https://goreportcard.com/badge/github.com/iseki0/goda)](https://goreportcard.com/report/github.com/iseki0/goda)
[![codecov](https://codecov.io/gh/iseki0/goda/graph/badge.svg?token=TBHUZUY561)](https://codecov.io/gh/iseki0/goda)

> **ThreeTen/JSR-310** 模型的 Go 实现

受 Java `java.time` 包（JSR-310）启发的 Go 实现，提供**类型安全**且**易于使用**的不可变日期和时间类型。

## 特性

### 核心类型

- 📅 **LocalDate**：不含时间的日期（例如：`2024-03-15`）
- ⏰ **LocalTime**：不含日期的时间（例如：`14:30:45.123456789`）
- 📆 **LocalDateTime**：不含时区的日期时间（例如：`2024-03-15T14:30:45.123456789`）
- 🌐 **ZoneOffset**：相对于格林威治/UTC 的时区偏移（例如：`+08:00`、`-05:00`、`Z`）
- 🌍 **OffsetDateTime**：带偏移的日期时间（例如：`2024-03-15T14:30:45.123456789+01:00`）
- 🔢 **Field**：日期时间字段枚举（类似 Java 的 `ChronoField`）
- 🔍 **TemporalAccessor**：用于查询时间对象的通用接口
- 📊 **TemporalValue**：带验证状态的类型安全字段值包装器

### 主要功能

- ✅ **ISO 8601 基本格式**支持（yyyy-MM-dd、HH:mm:ss[.nnnnnnnnn]，用 'T' 连接）
- ✅ **Java.time 兼容格式化**：小数秒对齐到 3 位数边界（毫秒、微秒、纳秒）
- ✅ **完整的 JSON 和 SQL** 数据库集成
- ✅ **日期运算**：支持溢出处理的天、月、年加减
- ✅ **类型安全的字段访问**：使用 `TemporalValue` 返回类型查询任何字段，验证支持和溢出
- ✅ **TemporalAccessor 接口**：跨所有时间类型的通用查询模式
- ✅ **链式操作**：流畅 API 配合错误处理进行复杂变更
- ✅ **不可变**：所有操作返回新值
- ✅ **类型安全**：通过不同类型实现编译时安全
- ✅ **零值友好**：正确处理零值

## 安装

```bash
go get github.com/iseki0/goda
```

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/iseki0/goda"
)

func main() {
    // 创建日期和时间
    date := goda.MustLocalDateOf(2024, goda.March, 15)
    time := goda.MustLocalTimeOf(14, 30, 45, 123456789)
    datetime := date.AtTime(time)  // 或 time.AtDate(date)

    fmt.Println(date)     // 2024-03-15
    fmt.Println(time)     // 14:30:45.123456789
    fmt.Println(datetime) // 2024-03-15T14:30:45.123456789

    // 直接从组件创建
    datetime2 := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 123456789)

    // 带时区偏移
    offset := goda.MustZoneOffsetOfHours(1)  // +01:00
    offsetDateTime := datetime.AtOffset(offset)
    fmt.Println(offsetDateTime) // 2024-03-15T14:30:45.123456789+01:00

    // 从字符串解析
    date, _ = goda.LocalDateParse("2024-03-15")
    time = goda.MustLocalTimeParse("14:30:45.123456789")
    datetime = goda.MustLocalDateTimeParse("2024-03-15T14:30:45")

    // 获取当前日期/时间
    today := goda.LocalDateNow()
    now := goda.LocalTimeNow()
    currentDateTime := goda.LocalDateTimeNow()
    currentOffsetDateTime := goda.OffsetDateTimeNow()

    // 日期运算
    tomorrow := today.Chain().PlusDays(1).MustGet()
    nextMonth := today.Chain().PlusMonths(1).MustGet()
    nextYear := today.Chain().PlusYears(1).MustGet()

    // 比较
    if tomorrow.IsAfter(today) {
        fmt.Println("明天在今天之后！")
    }
}
```

### 使用时区

```go
// 创建带偏移的时间
offset := goda.MustZoneOffsetOfHours(8)  // +08:00（中国标准时间）
odt := goda.MustNewOffsetDateTime(2024, goda.March, 15, 14, 30, 45, 0, offset)

// 解析带偏移的时间
odt, _ = goda.ParseOffsetDateTime("2024-03-15T14:30:45+08:00")
odt = goda.MustParseOffsetDateTime("2024-03-15T14:30:45Z")  // UTC

// 从 Go 的 time.Time 转换（保留偏移）
goTime := time.Now()
odt = goda.OffsetDateTimeOfGoTime(goTime)

// 改变偏移但保持本地时间
est := goda.MustZoneOffsetOfHours(-5)  // 美东时间
pst := goda.MustZoneOffsetOfHours(-8)  // 太平洋时间
odtEST := goda.MustNewOffsetDateTime(2024, goda.March, 15, 14, 30, 45, 0, est)
odtPST := odtEST.WithOffsetSameLocal(pst)  // 本地时间不变：14:30:45-08:00

// 改变偏移但保持瞬时时间
odtPST2 := odtEST.WithOffsetSameInstant(pst)  // 保持瞬时：11:30:45-08:00

// 带偏移的时间运算
tomorrow := odt.PlusDays(1)
inTwoHours := odt.PlusHours(2)

// 转换为 Unix 时间戳
epochSecond := odt.ToEpochSecond()

// 基于瞬时时间比较
if odt1.IsBefore(odt2) {
    fmt.Println("odt1 更早")
}
```

### 使用 TemporalValue 访问字段

使用 `Field` 枚举访问单个日期时间字段，返回类型安全的 `TemporalValue`：

```go
date := goda.MustLocalDateOf(2024, goda.March, 15)

// 检查字段支持
fmt.Println(date.IsSupportedField(goda.FieldDayOfMonth))  // true
fmt.Println(date.IsSupportedField(goda.FieldHourOfDay))   // false

// 获取带验证的字段值
year := date.GetField(goda.FieldYear)
if year.Valid() {
    fmt.Println("年份：", year.Int64())  // 2024
}

dayOfWeek := date.GetField(goda.FieldDayOfWeek)
if dayOfWeek.Valid() {
    fmt.Println("星期：", dayOfWeek.Int())  // 5（星期五）
}

// 不支持的字段返回 unsupported 的 TemporalValue
hourOfDay := date.GetField(goda.FieldHourOfDay)
if hourOfDay.Unsupported() {
    fmt.Println("LocalDate 不支持小时字段")
}

// 时间字段
time := goda.MustLocalTimeOf(14, 30, 45, 123456789)
hour := time.GetField(goda.FieldHourOfDay)
if hour.Valid() {
    fmt.Println("小时：", hour.Int())  // 14
}

nanoOfDay := time.GetField(goda.FieldNanoOfDay)
if nanoOfDay.Valid() {
    fmt.Println("自午夜以来的纳秒：", nanoOfDay.Int64())
}
```

**TemporalValue API：**
- `Valid() bool`：如果字段被支持且没有发生溢出，返回 true
- `Unsupported() bool`：如果该时间类型不支持该字段，返回 true
- `Overflow() bool`：如果字段值溢出，返回 true（保留供将来使用）
- `Int64() int64`：获取 int64 类型的原始值
- `Int() int`：获取 int 类型的值（为方便起见）

**为什么使用 TemporalValue？**

`TemporalValue` 返回类型提供类型安全的字段查询，防止静默错误：
- **明确的验证**：使用值之前检查 `Valid()`
- **清晰的错误语义**：区分不支持的字段和实际错误
- **面向未来**：准备好在需要时进行溢出检测
- **没有静默零值**：与原始 `int64` 返回不同，你可以区分 "0" 和 "不支持"

### TemporalAccessor 接口

所有时间类型都实现 `TemporalAccessor` 接口，提供统一的查询模式：

```go
// TemporalAccessor 提供对时间字段的只读访问
type TemporalAccessor interface {
    IsZero() bool
    IsSupportedField(field Field) bool
    GetField(field Field) TemporalValue
}

// 编写适用于任何时间类型的泛型函数
func printYear(t goda.TemporalAccessor) {
    if year := t.GetField(goda.YearField); year.Valid() {
        fmt.Printf("年份：%d\n", year.Int())
    }
}

// 适用于 LocalDate、LocalTime 或 LocalDateTime
printYear(goda.LocalDateNow())
printYear(goda.LocalDateTimeNow())
```

### 链式操作

所有时间类型都支持链式操作，用于流畅且带错误处理的复杂变更。链式操作允许你在单个表达式中执行多个修改，并进行适当的错误处理：

```go
// 流畅地链式多个操作
dt := goda.MustLocalDateTimeOf(2024, goda.March, 15, 14, 30, 45, 123456789)

// 链式日期和时间修改
meetingTime := dt.Chain().
    PlusDays(7).              // 下周
    WithHour(16).             // 下午 4 点
    WithMinute(0).            // 整点
    WithSecond(0).            // 无秒
    WithNano(0).              // 无纳秒
    MustGet()                 // 获取结果（出错时 panic）

fmt.Println("会议安排在：", meetingTime)

// 链式操作中的错误处理
result, err := dt.Chain().
    PlusMonths(1).
    WithDayOfMonth(32).       // 无效日期 - 会导致错误
    GetResult()               // 返回（零值，错误）

if err != nil {
    fmt.Println("无效操作：", err)
    // 使用后备方案
    validTime := dt.Chain().
        PlusMonths(1).
        WithDayOfMonth(31).   // 有效日期
        GetOrElse(dt)         // 出错时返回原始值
}
```

### JSON 序列化

```go
type Event struct {
    Name        string                `json:"name"`
    Date        goda.LocalDate        `json:"date"`
    Time        goda.LocalTime        `json:"time"`
    CreatedAt   goda.LocalDateTime    `json:"created_at"`
    ScheduledAt goda.OffsetDateTime   `json:"scheduled_at"`  // 带时区
}

event := Event{
    Name:        "会议",
    Date:        goda.MustLocalDateOf(2024, goda.March, 15),
    Time:        goda.MustLocalTimeOf(14, 30, 0, 0),
    CreatedAt:   goda.MustLocalDateTimeParse("2024-03-15T14:30:00"),
    ScheduledAt: goda.MustOffsetDateTimeParse("2024-03-15T14:30:00+08:00"),
}

jsonData, _ := json.Marshal(event)
// {"name":"会议","date":"2024-03-15","time":"14:30:00",
//  "created_at":"2024-03-15T14:30:00","scheduled_at":"2024-03-15T14:30:00+08:00"}
```

### 数据库集成

```go
type Record struct {
    ID          int64
    CreatedAt   goda.LocalDateTime
    Date        goda.LocalDate
    UpdatedAt   goda.OffsetDateTime  // 带时区用于审计日志
}

// 与 database/sql 一起工作 - 实现了 sql.Scanner 和 driver.Valuer
db.QueryRow("SELECT id, created_at, date, updated_at FROM records WHERE id = ?", 1).Scan(
    &record.ID, &record.CreatedAt, &record.Date, &record.UpdatedAt,
)

// 使用 offset datetime 插入
offset := goda.MustZoneOffsetOfHours(8)
now := goda.OffsetDateTimeNow()
db.Exec("INSERT INTO records (created_at, updated_at) VALUES (?, ?)",
    goda.LocalDateTimeNow(), now)
```

## API 概览

### 核心类型

| 类型                | 描述                                    | 示例                                   |
|---------------------|-----------------------------------------|----------------------------------------|
| `LocalDate`         | 不含时间的日期                          | `2024-03-15`                           |
| `LocalTime`         | 不含日期的时间                          | `14:30:45.123456789`                   |
| `LocalDateTime`     | 不含时区的日期时间                      | `2024-03-15T14:30:45`                  |
| `ZoneOffset`        | 相对于格林威治/UTC 的时区偏移           | `+08:00`、`-05:00`、`Z`                |
| `OffsetDateTime`    | 带 UTC 偏移的日期时间                   | `2024-03-15T14:30:45+08:00`            |
| `Month`             | 月份（1-12）                            | `March`                                |
| `Year`              | 年份                                    | `2024`                                 |
| `DayOfWeek`         | 星期（1=星期一，7=星期日）              | `Friday`                               |
| `Field`             | 日期时间字段枚举                        | `HourOfDay`、`DayOfMonth`              |
| `TemporalAccessor`  | 用于查询时间对象的接口                  | 所有时间类型都实现了此接口             |
| `TemporalValue`     | 带验证的类型安全字段值                  | 由 `GetField()` 返回                   |
| `Error`             | 带上下文的结构化错误                    | 提供详细的错误信息                     |
| `LocalDateChain`    | LocalDate 的链式操作                    | `date.Chain().PlusDays(1).MustGet()`   |
| `LocalTimeChain`    | LocalTime 的链式操作                    | `time.Chain().PlusHours(1).MustGet()`  |
| `LocalDateTimeChain`| LocalDateTime 的链式操作                | `dt.Chain().PlusDays(1).MustGet()`     |
| `OffsetDateTimeChain`| OffsetDateTime 的链式操作               | `odt.Chain().PlusHours(1).MustGet()`   |

### 格式规范

此包使用 ISO 8601 基本日历日期和时间格式（不是完整规范）：

**LocalDate**：`yyyy-MM-dd`（例如："2024-03-15"）  
仅限格里高利历日期。不支持周日期（YYYY-Www-D）或序数日期（YYYY-DDD）。

**LocalTime**：`HH:mm:ss[.nnnnnnnnn]`（例如："14:30:45.123456789"）  
24 小时格式。小数秒最多到纳秒。小数秒与 3 位数边界对齐（毫秒、微秒、纳秒），以实现 Java.time 兼容性：100ms → "14:30:45.100"，123.4ms → "14:30:45.123400"。解析接受任何长度的小数秒（例如："14:30:45.1" → 100ms）。

**LocalDateTime**：`yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]`（例如："2024-03-15T14:30:45.123456789"）  
使用 'T' 分隔符连接（解析时接受小写 't'）。

**ZoneOffset**：`±HH:mm[:ss]` 或 `Z` 表示 UTC（例如："+08:00"、" -05:30"、"Z"）  
小时数范围必须为 [-18, 18]，分钟和秒为 [0, 59]。还支持紧凑格式（±HH、±HHMM、±HHMMSS）。

**OffsetDateTime**：`yyyy-MM-ddTHH:mm:ss[.nnnnnnnnn]±HH:mm[:ss]`（例如："2024-03-15T14:30:45+08:00"）  
结合 LocalDateTime 和 ZoneOffset。接受 'Z' 作为 UTC 偏移。

### 时间格式化

时间值使用 ISO 8601 格式，并采用**与 Java.time 兼容**的小数秒对齐：

| 精度          | 位数 | 示例                                       |
|---------------|------|--------------------------------------------|
| 整秒          | 0    | `14:30:45`                                 |
| 毫秒          | 3    | `14:30:45.100`、`14:30:45.123`             |
| 微秒          | 6    | `14:30:45.123400`、`14:30:45.123456`       |
| 纳秒          | 9    | `14:30:45.000000001`、`14:30:45.123456789` |

小数秒自动对齐到 3 位数边界（毫秒、微秒、纳秒），与 Java 的 `LocalTime` 行为一致。解析接受任何长度的小数秒。

### 字段常量（30 个字段）

**时间字段**：`NanoOfSecond`、`NanoOfDay`、`MicroOfSecond`、`MicroOfDay`、`MilliOfSecond`、`MilliOfDay`、`SecondOfMinute`、`SecondOfDay`、`MinuteOfHour`、`MinuteOfDay`、`HourOfAmPm`、`ClockHourOfAmPm`、`HourOfDay`、`ClockHourOfDay`、`AmPmOfDay`

**日期字段**：`DayOfWeekField`、`DayOfMonth`、`DayOfYear`、`EpochDay`、`AlignedDayOfWeekInMonth`、`AlignedDayOfWeekInYear`、`AlignedWeekOfMonth`、`AlignedWeekOfYear`、`MonthOfYear`、`ProlepticMonth`、`YearOfEra`、`YearField`、`Era`

**其他字段**：`InstantSeconds`、`OffsetSeconds`

### 实现的接口

所有时间类型（`LocalDate`、`LocalTime`、`LocalDateTime`、`OffsetDateTime`）都实现了：
- `TemporalAccessor`：通用查询接口，使用 `GetField(field Field) TemporalValue`
- `fmt.Stringer`
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler`
- `encoding.TextAppender`
- `json.Marshaler` / `json.Unmarshaler`
- `sql.Scanner` / `driver.Valuer`

## 设计理念

此包遵循 **ThreeTen/JSR-310** 模型（Java 的 `java.time` 包），提供具有以下特点的日期和时间类型：

- **不可变**：所有操作返回新值
- **类型安全**：日期、时间和日期时间使用不同类型
- **简单格式**：使用 ISO 8601 基本格式（不是完整的复杂规范）
- **数据库友好**：直接集成 SQL
- **基于字段的访问**：通过 `TemporalAccessor` 接口的通用字段访问模式
- **安全的字段查询**：`TemporalValue` 返回类型验证字段支持并防止静默错误
- **零值安全**：在整个过程中正确处理零值

### 何时使用每种类型

**LocalDate、LocalTime、LocalDateTime** - 当时区不相关时使用：
- **生日**："3 月 15 日"在任何地方都表示 3 月 15 日
- **营业时间**：在本地上下文中的"上午 9:00 - 下午 5:00"
- **日程安排**：不考虑时区的"下午 2:30 会议"
- **日历日期**：历史日期、重复事件

**OffsetDateTime** - 当你需要相对于 UTC 的固定偏移时使用：
- **API 时间戳**：REST API 通常使用带偏移的 RFC3339
- **审计日志**：记录确切时刻及原始时区偏移
- **事件调度**：当时区偏移很重要但夏令时转换不重要时
- **国际协调**："会议在 UTC+1 的 14:00"

**ZoneOffset** - 用于表示时区偏移：
- **固定偏移**：+08:00、-05:00、Z（UTC）
- **不处理夏令时**：当不需要夏令时规则时使用
- **简单偏移运算**：在不同偏移之间转换

对于支持夏令时转换的完整时区支持，请使用 `ZonedDateTime`（即将推出）。

## 文档

完整的 API 文档可在 [pkg.go.dev](https://pkg.go.dev/github.com/iseki0/goda) 查看。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

此项目基于 MIT 许可证 - 详情请查看 LICENSE 文件。

