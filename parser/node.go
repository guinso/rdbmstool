package parser

const (
	NodeBool       NodeType = iota //a boolean constant
	NodeField                      //a db field such as table name, column name
	NodeIdentifier                 //a function name
	NodeNumber                     //a numerical constant
	NodeString                     //a string constant
	NodeParam                      //an SQL parameter
	NodeList                       //a list of SQL expression
	NodeCondition                  //a logical comparison expression; e.g. x > y
	NodeSelect                     //SQL select statement
	NodeFrom                       //SQL from statement
	NodeJoin                       //SQL join statement
	NodeWhere                      //SQL where statement
	NodeHaving                     //SQL having statement
	NodeGroupBy                    //SQL group by statement
	NodeLimit                      //SQl limit statement
	NodeUnion                      //SQL union statement
	NodeQuery                      //SQL query statement
)
