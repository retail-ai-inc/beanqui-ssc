<template>
  <div class="event">
    <div class="container-fluid">

        <div class="mb-3 row">

          <div class="col-2">
            <div class="row">
            <label class="col-sm-2 col-form-label">Id:</label>
            <div class="col-sm-10">
              <input type="text" class="form-control" id="staticEmail"  v-model="form.id">
            </div>
            </div>
          </div>

          <div class="col-2">
            <div class="row">
            <label class="col-sm-2 col-form-label">Status:</label>
            <div class="col-sm-10">
              <select class="form-select" aria-label="Default select" style="cursor: pointer" v-model="form.status">
                <option selected value="">Open this select</option>
                <option value="published">Published</option>
                <option value="success">Success</option>
                <option value="failed">Failed</option>
              </select>
            </div>
            </div>
          </div>

          <div class="col-2">
            <div class="col-auto">
              <button type="submit" class="btn btn-primary mb-3" @click="search">Search</button>
            </div>
          </div>

        </div>

      <table class="table">
        <thead>
          <tr>
            <th scope="col">#</th>
            <th scope="col">Id</th>
            <th scope="col">Channel</th>
            <th scope="col">Topic</th>
            <th scope="col">MoodType</th>
            <th scope="col">Status</th>
            <th scope="col">AddTime</th>
            <th scope="col">Payload</th>
            <th scope="col">Action</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, key) in eventLogs" :key="key">
            <th scope="row">{{parseInt(key)+1}}</th>
            <td>{{item.id}}</td>
            <td>{{item.channel}}</td>
            <td>{{item.topic}}</td>
            <td>{{item.moodType}}</td>
            <td>{{item.status}}</td>
            <td>{{item.addTime}}</td>
            <td>{{item.payload}}</td>
            <td>
              <div class="btn-group-sm" role="group">
                <button type="button" class="btn btn-primary dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
                  actions
                </button>
                <ul class="dropdown-menu">
                  <li><a class="dropdown-item" href="#">Retry</a></li>
                  <li><a class="dropdown-item" href="#">Delete</a></li>
                </ul>
              </div>
            </td>
          </tr>
        </tbody>

      </table>

    </div>
  </div>
</template>
<script setup>
import { reactive,onMounted,toRefs,onUnmounted } from "vue";
import request  from "request";

let data = reactive({
  eventLogs:[],
  page:1,
  pageSize:10,
  total:1,
  form:{
    id:"",
    status:""
  }
})

async function search(){
  let logs = await getEventLog(data.page,data.pageSize,data.form.id,data.form.status)
  data.eventLogs = logs.data
}

function getEventLog(page,pageSize,id,status){
  let params = {"page":page,"pageSize":pageSize,"id":"","status":""};
  if (id !== "") {
    params.id = id
  }
  if(status !== ""){
    params.status = status
  }
  return request.get("event_log/list",{"params":params});
}

onMounted(async()=>{
  let log =  await getEventLog(data.page,data.pageSize,data.form.id,data.form.status);
  data.eventLogs = log.data;
})

const {eventLogs,form} = toRefs(data);

</script>
<style scoped>
.event{
  transition: opacity 0.5s ease;
  opacity: 1;
}
</style>