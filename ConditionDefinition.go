package rdbmstool

//ConditionOperator logic operator for WHERE clause: =,<>,>,>=,<,<=, etc.
type ConditionOperator uint8

//Condition operator constants
const (
	//None no operator need to specify, use this couple with first condition
	None ConditionOperator = iota
	//And SQL AND condition operator
	And
	//Or SQL OR condition operator
	Or
)

func (operator ConditionOperator) String() string {
	switch operator {
	case And:
		return "AND"
	case Or:
		return "OR"
	default:
		return ""
	}
}

//ConditionDefinition SQL condition definition version 2
type ConditionDefinition struct {
	Condition        string
	ConditionComplex *ConditionDefinition
	Operator         ConditionOperator

	Conditions []ConditionDefinition
}

//NewCondition create a new binary condition
func NewCondition(expression string) *ConditionDefinition {
	return &ConditionDefinition{
		Condition:        expression,
		ConditionComplex: nil,
		Operator:         None,
		Conditions:       []ConditionDefinition{},
	}
}

//String generate SQL statement
func (cond *ConditionDefinition) String() (string, error) {
	sqlString := ""
	if cond.IsSimpleExpression() {
		sqlString = cond.Condition
	} else {
		tmpStr, tmpErr := cond.ConditionComplex.String()
		if tmpErr != nil {
			return "", tmpErr
		}

		sqlString = "(" + tmpStr + ")"
	}

	if len(cond.Conditions) == 0 {
		return sqlString, nil
	}

	for i := 0; i < len(cond.Conditions); i++ {
		tmpSQL, tmpErr := cond.Conditions[i].String()
		if tmpErr != nil {
			return "", tmpErr
		}

		sqlString += " " + cond.Conditions[i].Operator.String() + " " + tmpSQL
	}

	return sqlString, nil
}

//SetCondition set condition with expression string
func (cond *ConditionDefinition) SetCondition(condition string) *ConditionDefinition {
	cond.Condition = condition
	cond.Operator = None
	cond.ConditionComplex = nil
	cond.Conditions = nil

	return cond
}

//SetComplex set condition with ConditionDefinition instance
func (cond *ConditionDefinition) SetComplex(condDef *ConditionDefinition) *ConditionDefinition {
	cond.Condition = ""
	cond.Operator = None
	cond.ConditionComplex = condDef
	cond.Conditions = nil

	return cond
}

//AddAnd Append AND simple string expression condition
func (cond *ConditionDefinition) AddAnd(expression string) *ConditionDefinition {
	cond.Conditions = append(cond.Conditions, ConditionDefinition{
		Condition:        expression,
		ConditionComplex: nil,
		Operator:         And,
		Conditions:       nil,
	})

	return cond
}

//AddAndComplex Append AND nested condition
func (cond *ConditionDefinition) AddAndComplex(condition *ConditionDefinition) *ConditionDefinition {
	cond.Conditions = append(cond.Conditions, ConditionDefinition{
		Condition:        "",
		ConditionComplex: condition,
		Operator:         And,
		Conditions:       nil,
	})

	return cond
}

//AddOr Append AND simple string expression condition
func (cond *ConditionDefinition) AddOr(expression string) *ConditionDefinition {
	cond.Conditions = append(cond.Conditions, ConditionDefinition{
		Condition:        expression,
		ConditionComplex: nil,
		Operator:         Or,
		Conditions:       nil,
	})

	return cond
}

//AddOrComplex Append AND nested condition
func (cond *ConditionDefinition) AddOrComplex(condition *ConditionDefinition) *ConditionDefinition {
	cond.Conditions = append(cond.Conditions, ConditionDefinition{
		Condition:        "",
		ConditionComplex: condition,
		Operator:         Or,
		Conditions:       nil,
	})

	return cond
}

//AddComplex Append nested condition with ConditionDefinition instance
func (cond *ConditionDefinition) AddComplex(
	operator ConditionOperator, condition *ConditionDefinition) *ConditionDefinition {

	cond.Conditions = append(cond.Conditions, ConditionDefinition{
		Condition:        "",
		ConditionComplex: condition,
		Operator:         operator,
		Conditions:       nil,
	})

	return cond
}

//GetConditions get condition by index
func (cond *ConditionDefinition) GetConditions(index int) *ConditionDefinition {
	if index >= 0 && index < len(cond.Conditions) {
		return &cond.Conditions[index]
	}

	return nil
}

//IsSimpleExpression check first condition is simple expression instead of nested condition
func (cond *ConditionDefinition) IsSimpleExpression() bool {
	return cond.ConditionComplex == nil
}
