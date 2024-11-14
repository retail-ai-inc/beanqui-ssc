const eventApi = {
    // make bootstrap alert html element
    Alert(message,type){
        const alertPlaceholder = document.getElementById('payloadAlertInfo');
        alertPlaceholder.innerHTML = `<div class="alert alert-${type} alert-dismissible" id="my-alert" role="alert">
      <div>${message}</div>
      <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
      </div>`;
    },
    Delete(id){
        let params = {id:id};
        return request.delete(`event_log/delete`,{params});
    },
    Edit(id,payload){
        const headers = {
            "Content-Type":"application/x-www-form-urlencoded"
        }
        return request.put(`/event_log/edit`,{id:id,payload:payload},{headers:headers});
    },
    Retry(id,data){
        return request.post(`/event_log/retry`,{uniqueId:id,data:JSON.stringify(data)});
    }
}