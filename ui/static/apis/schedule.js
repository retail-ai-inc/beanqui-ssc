const schedule = {
    GetSchedule(page,pageSize){
        return request.get("schedule",{"params":{"page":page,"pageSize":pageSize}});
    }
}