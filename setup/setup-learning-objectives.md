# Setup learning objectives steps 

ok - create grading class 
    GradingSystem
        Grade 
            Highest
            Lowest

ok - add grading to activity

## How to add grading assignment the code version that was used?
    - ok Assignment
    - ok Grade value - 0 to 100
    - ok skill
    - CodeHistory

## Code versioning
    - Add CodeHistory
    - Add version to Code that should be auto updated  
    ds.DB.Create(&luchador) and update
    
    add hooks to Code 

    - Create Code hook on add/update actions to create records in CodeHistory
        make sure there is no concurrency in this action
        routes\play\request_handler.go: 39 use mutex

    - Add CodeHistory reference to AssignmentGrade
    - When creating AssignmentGrade link to CodeHistory

** Setup process 

- load & add if dont exist learning objective
- load & add if dont exist skills h
- load & add if dont exist grading
- load activity list with skill placeholder objects containing name attribute only
- populate skills objects based on the name in activity list
- add if dont exist activities

create 3 folders under learning-definition folder
- learning-experience
- skill
- activity

Inside each folder add json files with the records to be created.
Add extra parameter to the docker file with the learning-definition folder
