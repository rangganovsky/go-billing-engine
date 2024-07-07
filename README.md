This is Test Project to simulate billing process, we will build it using golang with GIN, postgresql (with GORM),

the engine should provide ability to 
- Loan schedule for a given loan( when am i supposed to pay how much)
- Outstanding Amount for a given loan
- Status of weather the borrower is Delinquent or not

we also need to provide seeding the data for the Borrower and Loan, thus we provide simple CRUD for both entities

Scenarios : 
- Borrower repay the amount every week ( assume they can only pay the exact amount payable that week or not pay at all)
- Borrower can be flag as delinquent if they miss 2 continuous repayment

TODO : 
- add unit test
- make better error handling
- create event-driven case to track loan and billing schedule updated status. now we are still using the simples method which is ON DEMAND
- provide data seeder