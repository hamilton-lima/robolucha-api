# Setup learning objectives steps 

- create grading class 
    GradingSystem
        Grade 
            Highest
            Lowest
            
- add grading to activity

- load & add if dont exist learning experience
- load & add if dont exist skills 
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
