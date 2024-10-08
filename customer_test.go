package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

type CustomerTestSteps struct {
	customerService *CustomerService
	err             error
	count           int
}

var DEFAULT_BIRTHDAY = time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func (t *CustomerTestSteps) theCustomerIsCreated(ctx context.Context, fn, ln string) error {
	t.err = t.customerService.AddCustomer(fn, ln, DEFAULT_BIRTHDAY)
	return nil
}

func (t *CustomerTestSteps) theCustomerCreationShouldBeSuccessful(ctx context.Context) error {
	if t.err != nil {
		return fmt.Errorf("expected no error but got %v", t.err)
	}
	return nil
}

func (t *CustomerTestSteps) theCustomerCreationShouldFail(ctx context.Context) error {
	if t.err == nil {
		return fmt.Errorf("expected error but got nil")
	}

	if t.err.Error() != "mandatory name parameter is missing" {
		return fmt.Errorf("expected 'mandatory name parameter is missing' error but got '%s'", t.err.Error())
	}

	return nil
}

func (t *CustomerTestSteps) theSecondCustomerCreationShouldFail(ctx context.Context) error {
	if t.err == nil {
		return fmt.Errorf("expected error but got nil")
	}

	if t.err.Error() != "customer already exists" {
		return fmt.Errorf("expected 'customer already exists' error but got '%s'", t.err.Error())
	}

	return nil
}

func (t *CustomerTestSteps) thereAreSomeCustomers(ctx context.Context, table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 {
			continue // skip header...
		}

		t.customerService.AddCustomer(row.Cells[0].Value, row.Cells[1].Value, DEFAULT_BIRTHDAY)
	}
	return nil
}

func (t *CustomerTestSteps) allCustomersAreSearched(ctx context.Context) error {
	t.count = len(t.customerService.SearchCustomers())
	return nil
}

func (t *CustomerTestSteps) theCustomerSabineMustermannIsSearched(ctx context.Context, fn, ln string) error {
	t.count = len(t.customerService.SearchCustomersByName(fn, ln))
	return nil
}

func (t *CustomerTestSteps) theCustomerFnLnCanBeFound(ctx context.Context, fn, ln string) error {
	customer := t.customerService.SearchCustomer(fn, ln)

	if customer.FirstName != fn {
		return fmt.Errorf("expected first name to be Sabine but got %v", customer.FirstName)
	}

	if customer.LastName != ln {
		return fmt.Errorf("expected last name to be Mustermann but got %v", customer.LastName)
	}

	return nil
}

func (t *CustomerTestSteps) theNumberOfCustomersFoundIs(ctx context.Context, expectedCount int) error {
	if t.count != expectedCount {
		return fmt.Errorf("expected %d customers to be found but got %d", expectedCount, t.count)
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {

	t := CustomerTestSteps{
		customerService: NewCustomerService(),
		err:             nil,
		count:           0,
	}

	sc.Given(`^the customer (\w+) (\w+) is created$`, t.theCustomerIsCreated)
	sc.When(`^the customer (\w+) (\w+) is created$`, t.theCustomerIsCreated)
	sc.When(`an invalid customer (\w*) (\w*) is created`, t.theCustomerIsCreated)
	sc.When(`the second customer (\w+) (\w+) is created`, t.theCustomerIsCreated)
	sc.Then(`the customer creation should be successful`, t.theCustomerCreationShouldBeSuccessful)
	sc.Then(`the customer creation should fail`, t.theCustomerCreationShouldFail)
	sc.Then(`the second customer creation should fail`, t.theSecondCustomerCreationShouldFail)
	sc.Given(`there is a customer`, t.thereAreSomeCustomers)
	sc.Given(`there are some customers`, t.thereAreSomeCustomers)
	sc.When(`all customers are searched`, t.allCustomersAreSearched)
	sc.When(`the customer (\w+) (\w+) is searched`, t.theCustomerSabineMustermannIsSearched)
	sc.Then(`the customer (\w+) (\w+) can be found`, t.theCustomerFnLnCanBeFound)
	sc.Then(`^the number of customers found is (\d+)$`, t.theNumberOfCustomersFoundIs)

}
