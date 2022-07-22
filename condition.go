// OpenCypher9 Spec - Types of constraints
// ---------------------------------------
// C - Completed
// IP - In progress
// P - Planned
// NS - Not supported
// ---------------------------------------
// [C] 1. Boolean operations (n.name = 'Peter' XOR (n.age < 30 AND n.name = 'Tobias') OR NOT (n.name = 'Tobias' OR n.name = 'Peter'))
// [C] 2. Filter on node label (n:Swedish)
// [C] 3. Filter on node property (n.age > 30)
// [?] 4. Filter on relationship property (k.since < 2000)
// [?] 5. Filter on dynamically-computed node property (n[toLower(propname)]< 30)
// [C] 6. Property existence checking (exists(n.belt))
// [C] 7. Match the beginning of a string (n.name STARTS WITH 'Pet')
// [C] 8. Match the ending of a string (n.name ENDS WITH 'ter')
// [C] 9. Match anywhere within a string (n.name CONTAINS 'ete')
// [C] 10. String matching negation (NOT n.name ENDS WITH 's')
// [C] 11. Filter on patterns (others.name IN ['Andres', 'Peter'] AND (tobias)<-[]-(others))
// [C] 12. Filter on patterns using NOT (NOT (persons)-[]->(peter))
// [C] 13. Filter on patterns with properties ((n)-[:KNOWS]-({name: 'Tobias'}))
// [?] 14. Filter on relationship type (n.name='Andres' AND type(r) STARTS WITH 'K')
// [C] 15. IN operator (a.name IN ['Peter', 'Tobias'])
// [?] 16. Default to false if property is missing (n.belt = 'white')
// [?] 17. Default to true if property is missing (n.belt = 'white' OR n.belt IS NULL RETURN n.name, n.age, n.belt)
// [?] 18. Filter on null (person.name = 'Peter' AND person.belt IS NULL RETURN person.name, person.age, person.belt)
// [C] 19. Simple range (a.name >= 'Peter')
// [C] 20. Composite range (a.name > 'Andres' AND a.name < 'Tobias')

package go_cypherdsl

import (
	"errors"
	"fmt"
	"strings"
)

type ConditionBuilder struct {
	Start   *operatorNode
	Current *conditionNode
	errors  []error
}

// C is used to start a condition chain provided with a ConditionConfig
func C(condition *ConditionConfig) ConditionOperator {
	wq, err := NewCondition(condition)

	if err != nil {
		return &ConditionBuilder{
			errors: []error{err},
		}
	}

	cond := &ConditionBuilder{
		Start:   nil,
		errors:  nil,
		Current: nil,
	}

	node := &operatorNode{
		First: true,
		Query: &conditionNode{
			Condition: wq,
			Next:      nil,
		},
	}

	err = cond.addNext(node)
	if err != nil {
		cond.addError(err)
	}

	return cond
}

//add an error to the condition chain
func (c *ConditionBuilder) addError(e error) {
	if c.errors == nil {
		c.errors = []error{e}
	} else {
		c.errors = append(c.errors, e)
	}
}

//check if the builder has had any errors down the chain
func (c *ConditionBuilder) hasErrors() bool {
	return c.errors != nil && len(c.errors) > 0
}

//add another node to the chain
func (c *ConditionBuilder) addNext(node *operatorNode) error {
	if node == nil {
		return errors.New("node can not be nil")
	}

	if node.Query == nil {
		return errors.New("next can not be nil")
	}

	//different behavior if its the first of the chain
	if c.Start == nil {
		c.Start = node
		c.Current = node.Query
	} else {
		c.Current.Next = node
		c.Current = node.Query
	}

	return nil
}

func (c *ConditionBuilder) And(condition *ConditionConfig) ConditionOperator {
	return c.addCondition(condition, "AND")
}

func (c *ConditionBuilder) Or(condition *ConditionConfig) ConditionOperator {
	return c.addCondition(condition, "OR")
}

func (c *ConditionBuilder) Xor(condition *ConditionConfig) ConditionOperator {
	return c.addCondition(condition, "XOR")
}

func (c *ConditionBuilder) Not(condition *ConditionConfig) ConditionOperator {
	return c.addCondition(condition, "NOT")
}

func (c *ConditionBuilder) AndNested(query WhereQuery, err error) ConditionOperator {
	return c.addNestedCondition(query, err, "AND")
}

func (c *ConditionBuilder) OrNested(query WhereQuery, err error) ConditionOperator {
	return c.addNestedCondition(query, err, "OR")
}

func (c *ConditionBuilder) XorNested(query WhereQuery, err error) ConditionOperator {
	return c.addNestedCondition(query, err, "XOR")
}

func (c *ConditionBuilder) NotNested(query WhereQuery, err error) ConditionOperator {
	return c.addNestedCondition(query, err, "NOT")
}

func (c *ConditionBuilder) addNestedCondition(query WhereQuery, err error, condType string) ConditionOperator {
	if c.hasErrors() {
		return c
	}

	if err != nil {
		c.addError(err)
		return c
	}

	//create node, make sure to wrap the query in parenthases since its nested.
	node := &operatorNode{
		Condition: condType,
		Query: &conditionNode{
			Condition: "(" + query + ")",
			Next:      nil,
		},
	}

	//add it
	err = c.addNext(node)
	if err != nil {
		c.addError(err)
	}

	//return pointer to the builder
	return c
}

func (c *ConditionBuilder) addCondition(condition *ConditionConfig, condType string) ConditionOperator {
	//check if any errors are present, if they are, bail
	if c.hasErrors() {
		return c
	}

	//convert condition object into actual cypher
	wq, err := NewCondition(condition)
	if err != nil {
		c.addError(err)
		return c
	}

	//create the next node of the linked list
	node := &operatorNode{
		Condition: condType,
		Query: &conditionNode{
			Condition: wq,
			Next:      nil,
		},
	}

	//add it
	err = c.addNext(node)
	if err != nil {
		c.addError(err)
	}

	//return pointer to the builder
	return c
}

func (c *ConditionBuilder) Build() (WhereQuery, error) {
	//if it has errors, compile that and return
	if c.hasErrors() {
		errStr := ""
		for _, err := range c.errors {
			errStr += err.Error() + ";"
		}

		errStr = strings.TrimSuffix(errStr, ";")

		return "", fmt.Errorf("(%v) errors occurred: %s", len(c.errors), errStr)
	}

	query := ""

	//if start is not defined, something went wrong
	if c.Start == nil {
		return "", errors.New("no condition defined")
	}

	i := c.Start

	//iterate...
	for {
		if i == nil || i.Query == nil {
			break
		}

		t := ""

		if i.First {

		} else {
			t += i.Condition + " "
		}

		query += t + i.Query.Condition.ToString() + " "

		//iterate up
		i = i.Query.Next
	}

	//return entire condition
	return WhereQuery(strings.TrimSuffix(query, " ")), nil
}

type conditionNode struct {
	Condition WhereQuery
	Next      *operatorNode
}

type operatorNode struct {
	First     bool
	Condition string
	Query     *conditionNode
}

type ConditionOperator interface {
	And(c *ConditionConfig) ConditionOperator
	AndNested(query WhereQuery, err error) ConditionOperator
	Or(c *ConditionConfig) ConditionOperator
	OrNested(query WhereQuery, err error) ConditionOperator
	Xor(c *ConditionConfig) ConditionOperator
	XorNested(query WhereQuery, err error) ConditionOperator
	Not(c *ConditionConfig) ConditionOperator
	NotNested(query WhereQuery, err error) ConditionOperator
	Build() (WhereQuery, error)
}

type BooleanOperator string

const (
	LessThanOperator             BooleanOperator = "<"
	GreaterThanOperator          BooleanOperator = ">"
	LessThanOrEqualToOperator    BooleanOperator = "<="
	GreaterThanOrEqualToOperator BooleanOperator = ">="
	EqualToOperator              BooleanOperator = "="  // TODO: Rename to `EqualityOperator` to match OC9 spec
	NotEqualToOperator           BooleanOperator = "<>" // TODO: Rename to `InequalityOperator` to match OC9 spec
	InOperator                   BooleanOperator = "IN"
	IsOperator                   BooleanOperator = "IS"
	RegexEqualToOperator         BooleanOperator = "=~"
	StartsWithOperator           BooleanOperator = "STARTS WITH"
	EndsWithOperator             BooleanOperator = "ENDS WITH"
	ContainsOperator             BooleanOperator = "CONTAINS"
)

func (b BooleanOperator) String() string {
	return string(b)
}

// ConditionConfig is the configuration object for where conditions
type ConditionConfig struct {
	// Either ConditionOperator or ConditionFunction can be set, depending on whether an operator (e.g. '=') is
	// required, or a function (e.g. exists()). Both cannot be specified.
	ConditionOperator BooleanOperator

	//condition functions that can be used
	ConditionFunction string

	// When constructing a conditional, Name must be specified (unless filtering on a Path). Additionally, either
	// Field or Label must be specified (but not both). If Field is specified, FieldManipulationFunction is an
	// optional function that can be specified to manipulate the specified Field during the query. If Label
	// is specified, all other fields (other than Name and NegateCondition) will be ignored.
	Name  string
	Field string // TODO: Rename to `Property` to match OC9 spec
	Label string

	//exclude parentheses
	FieldManipulationFunction string

	// When filtering on a Path, the only other field that can be specified is NegateCondition.
	Path *PathBuilder

	// When using any operator to compare to a specific value (other than InOperator), this field must be specified.
	// If Check and (CheckName, CheckField) are both specified - (CheckName, CheckField) will take precedence.
	Check interface{}

	// When comparing one node to another, CheckField and CheckName must be specified.
	CheckName  string
	CheckField string

	// When using the InOperator, this field must be specified. This will not apply to any other operators.
	CheckSlice []interface{}

	// When NegateCondition is set to true, NOT is appended to the start of this condition to apply a logical
	// negation to the statement.
	NegateCondition bool
}

func (condition *ConditionConfig) ToString() (string, error) {
	var sb strings.Builder

	if condition.NegateCondition {
		sb.WriteString("NOT ")
	}

	if condition.Path != nil {
		cypher, err := condition.Path.ToCypher()
		if err != nil {
			return "", fmt.Errorf("error converting path to cypher: %s", err)
		}
		sb.WriteString(cypher)
		return sb.String(), nil
	}

	//check initial error conditions
	if condition.Name == "" {
		return "", errors.New("to construct a query, Name must be specified")
	}

	if condition.Label == "" && (condition.Field == "" && condition.FieldManipulationFunction == "") {
		return "", errors.New("either Label or (Field, FieldManipulationFunction) must be specified")
	}

	if condition.Field != "" && condition.Label != "" {
		return "", errors.New("only one of Field and Label can be specified")
	}

	if condition.FieldManipulationFunction != "" && condition.Field == "" {
		return "", errors.New("if FieldManipulationFunction is specified, Field must also be specified")
	}

	if condition.Label != "" {
		sb.WriteString(condition.Name)
		sb.WriteRune(':')
		sb.WriteString(condition.Label)
		return sb.String(), nil
	}

	var node string
	if condition.Field != "" {
		node = fmt.Sprintf("%s.%s", condition.Name, condition.Field)
	} else {
		node = condition.Name
	}

	if condition.FieldManipulationFunction != "" {
		sb.WriteString(condition.FieldManipulationFunction)
		sb.WriteRune('(')
		sb.WriteString(node)
		sb.WriteRune(')')
	} else {
		sb.WriteString(node)
	}

	if condition.ConditionOperator == "" && condition.ConditionFunction == "" {
		return "", errors.New("one of ConditionOperator or ConditionFunction must be specified")
	}

	if condition.ConditionOperator != "" && condition.ConditionFunction != "" {
		return "", errors.New("only one of ConditionOperator or ConditionFunction can be specified")
	}

	// Handle ConditionFunction if specified
	if condition.ConditionFunction != "" {
		if condition.NegateCondition {
			return fmt.Sprintf(
				"NOT %s(%s)", condition.ConditionFunction, strings.Trim(sb.String(), "NOT "),
			), nil
		}
		return fmt.Sprintf("%s(%s)", condition.ConditionFunction, sb.String()), nil
	}

	// Add ConditionOperator to query
	sb.WriteRune(' ')
	sb.WriteString(condition.ConditionOperator.String())
	// Handle edge cases for specific ConditionOperator types
	if condition.ConditionOperator == InOperator {
		if condition.CheckSlice == nil {
			return "", errors.New("slice can not be nil")
		}

		if condition.Check != nil {
			return "", errors.New("check should not be defined when using in operator")
		}

		if len(condition.CheckSlice) == 0 {
			return "", errors.New("slice should not be nil")
		}

		q := "["

		for _, val := range condition.CheckSlice {
			str, err := cypherizeInterface(val)
			if err != nil {
				return "", err
			}

			q += fmt.Sprintf("%s,", str)
		}

		sb.WriteRune(' ')
		sb.WriteString(strings.TrimSuffix(q, ","))
		sb.WriteRune(']')
	} else {
		if condition.CheckName != "" && condition.CheckField != "" {
			sb.WriteRune(' ')
			sb.WriteString(condition.CheckName)
			sb.WriteRune('.')
			sb.WriteString(condition.CheckField)
		} else {
			str, err := cypherizeInterface(condition.Check)
			if err != nil {
				return "", err
			}
			sb.WriteRune(' ')
			sb.WriteString(str)
		}
	}

	return sb.String(), nil
}

func NewCondition(condition *ConditionConfig) (WhereQuery, error) {
	if condition == nil {
		return "", errors.New("condition can not be nil")
	}

	str, err := condition.ToString()
	if err != nil {
		return "", err
	}

	return WhereQuery(str), nil
}
