1. Using "Command prompts" (links to the prompts in the description)

2. Buiding features:

    2.1 plan_feature command.

    2.2 review the plan

    2.3 phases

    2.4 code_review command

    2.5 Read through review

    2.6 Select fixes to implement.
    
    2.7 Manually test.

3. The code structure folders:

    3.1 `cmd`: This is the folder contains all the handlers + initial service needed.

    3.2 `task_docs`: This is generated docs for cursor.

    3.3 `internal/components`: Here contains all the handlers (Logic will sit here)

    3.4 `internal/pkg`: Contains all the package necessary for application

    3.5 `internal/repositories`: Contains all the database table schema and some CRUD for each table there. Each file will have individual table.

    3.6 `internal/services`: Contains all the code for calling external service
