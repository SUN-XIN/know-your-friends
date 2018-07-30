# Time Organization  
* Phase1: understand subject (no anything about technique)  
Understand the objective, especially the definition of categories: 
    Best Friend, Crush, Most seen and Mutual Love.
Imagine all the use cases, and think about the exceptional case for each category.

* Phase2: basic thinking about technique  
Tool to use (Protobuf + gRPC + Kafka)
Input/Output
DB kind (Redis + ScyllaDB)
known/unknown element

* Phase3: explore the unknown element (Protobuf + gRPC)  

* Phase4: explore the unknown element (ScyllaDB)  

* Phase5: pseudo code (frame) + find out some difficulties   

* Phase6: resolve the difficulties and fill the code  

* Phase7: finish the code and test with the different work-flows  

* Phase8: Doc

# Analyze subject  
* relation of categories  
Most seen: most duration  
Best Friend: "Most seen" + outside of significant place  
Mutual Love: "Most seen" + each other  
Crush: in home (geo) + in night (time) + count (last rolling N days)  

* property of categories  
mutual: if A is in a category of B, so B must also be in in the same category of A ->  
Crush, Mutual Love  
independent ->    
Most seen, Best Friend  

# DB architecture   
ScyllaDB, nosql, can not do the complicated query.  
No transaction.  

# Cache and checkpoint
Local Cache: manage some constraints, such like DAY, size.   
    if not found in local cache, get it from DB.  
Save the result of day in DB, we can then use it as checkpoint when server restarts

# Work-flow  
1.  Server receives 1 session {user1,user2,startDate,endDate,lat,lng}  
2.  Check if user1's places is in local cache (get from DB if not)  
3.  Check if this session is in any of user1's places, save the result in `isIn`  
4.  Put in the table SessionIntegrate whose PRIMARY KEY is `use_id_1 + use_id_2 + day`  
5.  Check if there is previous result for user1 in local cache    
6.  check MostSeen anyway, and if isIn is false, check BestFriend     
7.  not in cahce + not in DB -> first session of the day    
    Create the result in DB and put it in local cache  
8.  in local cache (or find in db)   
    if user1's MostSeen found (is the same day of this session) + (is user2) ->  update cache, update result in DB  
    if user1's MostSeen found (is the same day of this session) + (is not user2) -> 
    recalculate the result by fetching data from DB, then update cache, update result in DB  
    if user1's MostSeen found (is not the same day of this session)  -> 
    recalculate the result by fetching data from DB, then update cache, update result in DB  
9.  Check Crush  
    in home (either user1 or user2) + session is more than 6h + in night -> put/update in the table session_crush, whose PRIMARY KEY is `(user_id_owner + day)`  
    then check NB sessions in the last 7 days  

# Test client (/client_simple/main.go)
* CheckInitData  
when DB is empty, init some sessions  
* CheckWorkflowCrush  
create 3 sessions which are: in place (home) + at night + but duration is less than 6h ->   
only MostSeen in the response  
create 3 sessions which are: in place (home) + duration 6h30, but not at night ->   
only MostSeen in the response  
create 3 sessions which are: in place (home) + at night + duration 6h30 ->   
MostSeen and Crush in the response  
create 3 sessions which are: another friendID + in place (home) + at night + duration 6h30 ->   
MostSeen and Crush(2 diffenrent friend IDs) in the response  
* CheckWorkflowMostSeen
create 1 session which is: int places, not in night ->  
only MostSeen in the response   
create 1 session which is: another friendID, more duration + another day + in places + not in night ->  
MostSeen with the new friendID in the response  
create 1 session which is: the first friendID, much more duration + another day + in places + not in night ->  
MostSeen with the first friendID in the response 
* CheckWorkflowBestFriends
create 1 session which is: out places, not in night ->  
MostSeen/BestFriend in the response   
create 1 session which is: another friend, more duration + another day + out places + not in night ->  
MostSeen/BestFriend with the new friendID in the response  
create 1 session which is: friends, more duration + another day + out places + not in night ->  
MostSeen/BestFriend with the first friendID in the response 
* CheckMutualLove  
1/  owner spends 100s with friend1   
2/  friend2 spends 300s with friend3  
3.1/ owner spends 200s with friend2  
3.2/ friend2 spends 200s with owner  
4.1/ owner spends 200s with friend2  
4.2/ friend2 spends 200s with owner  

# Optimisation/TODO  
* time zone  
    UTC or user's time zone  
* the session through 2 days  
    This DB architecture use DAY as time unit, there will be the session through 2 days.  
    ex: StartDate is 2018-07-28 23h00 EndDate is 2018-07-29 01h00  
    Use EndDate as the DAY of DB  
* Super long session (more than 24h)  
    Need to cut into sub-sessions  
* Few data  
    ex: when the service starts, no any data. So 0 session inthe last rolling 7 days.  
* problem of neighbourhood (2D-3D)  
    geo is 2 dimensions x and y, but not 3d (height dimension is missing)  
    ex: A lives in stage 2, and B is his neighbourhood in stage 3. Lat/Lng is the same for A and B.  
* past session  
    This service need to process past session?  
    ex: session of last year/month/week  
    when we need to process 1 session of 2018-01-07, it will change the result for days: 2018-01-08, 2018-01-09, 2018-01-10 ...
* Crush calculate only one time
    if A is Crush of B, so B is also Crush of A
* retry (TECH)  
    retry when failed  
    ex: when scylladb failed  
* validate data (TECH)  
    validate data before insert in db  
    validate data when receive session  
    etc.  
* local cache for Crush (TECH)  
    now always get sessions from db and recalculate for each new session
* heartbeats/stat endpoint (TECH)  
    heartbeats: if server is running  
    stat: what happende in server (monitoring)  
* alert when down/restart  
* configuration   
    last rolling N days  
    N nights for "Crush"  
    etc.  