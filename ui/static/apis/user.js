const userApi = {
    List(){
        return request.get("/user/list");
    },
    Add(data){
        return request.post("/user/add",data);
    },
    Delete(account){
        let params = {account:account};
        return request.post(`/user/del`,params,{headers:{"Content-Type":"application/json"}});
    },
    Edit(data){
        const headers = {
            "Content-Type":"application/x-www-form-urlencoded"
        }
        return request.post(`/user/edit`,data,{headers:headers});
    }
}