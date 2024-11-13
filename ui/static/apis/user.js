const userApi = {
    List(){
        return request.get("/user/list");
    },
    Add(data){
        return request.post("/user/add",data);
    },
    Delete(account){
        return request.delete(`/user/del?account=${account}`);
    },
    Edit(data){
        const headers = {
            "Content-Type":"application/x-www-form-urlencoded"
        }
        return request.put(`/user/edit`,data,{headers:headers});
    }
}