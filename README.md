# rdbmstool
RDMS helper tool for Golang project

SQL syntax is compatibles with SQL:2003 standard to allow same SQL statement applicable across major database vendor

# Example
```golang
sqlQuery := rdbmstool.NewQueryBuilder()

sqlQuery.From("role_access", "a").
    Select("a.id", "").
    Select("a.role_id", "").
    Select("a.access_id", "").
    Select("b.name", "role").
    Select("c.name", "access").
    Select("a.is_authorize", "").
    JoinSimple("role", "b", rdbmstool.LEFT_JOIN, "a.role_id", "b.id", rdbmstool.EQUAL).
    JoinSimple("access", "c", rdbmstool.LEFT_JOIN, "a.access_id", "c.id", rdbmstool.EQUAL).
    //first where statement can use either 'WhereOR()' or 'WhereAnd' (they are no difference)
    WhereOR(rdbmstool.LIKE, "b.name", "'%keyword%'"). 
    WhereOR(rdbmstool.LIKE, "c.name", "'%asdqwe%'")
    Limit(10, 2) //page index (value 2) is zero based
    OrderBy("a.id", true)
    OrderBy("access", false)
    GroupBy("a.role_id", false)

sqlStr := sqlQuery.SQL() //output into SQL statement string
```
**SQL output**
```sql
SELECT a.id, a.role_id, a.access_id, b.name AS role, c.name AS access, a.is_authorize
FROM role_access a
LEFT JOIN role b ON a.role_id = b.id
LEFT JOIN access c ON a.access_id = c.id
WHERE b.name LIKE '%keyword%' OR c.name LIKE '%asdqwe%'
GROUP BY a.role_id DESC
ORDER BY a.id ASC, access DESC
LIMIT 10 OFFSET 20
```