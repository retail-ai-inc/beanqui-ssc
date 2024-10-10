<template>
  <div class="event">

    <div class="container-fluid">

        <div class="mb-3 row">

          <div class="col-2">
            <div class="row">
            <label for="formId" class="col-sm-2 col-form-label">Id:</label>
            <div class="col-sm-10">
              <input type="text" class="form-control" id="formId" name="formId"  v-model="form.id">
            </div>
            </div>
          </div>

          <div class="col-2">
            <div class="row">
              <label for="formStatus" class="col-sm-2 col-form-label">Status:</label>
              <div class="col-sm-10">
                <select class="form-select" aria-label="Default select" id="formStatus" name="formStatus" style="cursor: pointer" v-model="form.status">
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
      <Pagination :page="page" :total="total" :cursor="cursor" @changePage="changePage"/>
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
          <tr v-for="(item, key) in eventLogs" :key="key" style="height: 3rem;line-height:3rem">
            <th scope="row">{{parseInt(key)+1}}</th>
            <td><router-link to="" class="nav-link text-muted" style="display: contents" v-on:click="detailEvent(item)">{{item.id}}</router-link></td>
            <td>{{item.channel}}</td>
            <td>{{item.topic}}</td>
            <td>{{item.moodType}}</td>
            <td>
              <span v-if="item.status == 'success'" class="text-success">{{item.status}}</span>
              <span v-else-if="item.status == 'failed'" class="text-danger">{{item.status}}</span>
              <span v-else-if="item.status == 'published'" class="text-warning">{{item.status}}</span>
            </td>
            <td>{{item.addTime}}</td>
            <td>
              <span class="d-block text-truncate" style="max-width: 30rem;">
                {{item.payload}}
              </span>
            </td>
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
      <Pagination :page="page" :total="total" :cursor="cursor" @changePage="changePage"/>
    </div>
  </div>
</template>
<script setup>
import { reactive,onMounted,toRefs,onUnmounted } from "vue";
import { useRouter } from 'vueRouter';
import request  from "request";
import Pagination from "../components/pagination.vue";
import cfg  from "config";

let data = reactive({
  eventLogs:[],
  page:1,
  pageSize:10,
  total:1,
  cursor:0,
  form:{
    id:"",
    status:""
  },
  sseEvent:null
})

async function search(){

  sessionStorage.setItem("id",data.form.id);
  sessionStorage.setItem("status",data.form.status);

  initEventSource();
}

async function changePage(page,cursor){
  data.page = page;
  data.cursor = cursor;
  sessionStorage.setItem("page",page)

  initEventSource();
}
const uRouter = useRouter();
function detailEvent(item){
  uRouter.push("detail/"+item.id);
}

function initEventSource(){
  let url = `${cfg.sseUrl}event_log/list?page=${data.page}&pageSize=${data.pageSize}&id=${data.form.id}&status=${data.form.status}&token=${sessionStorage.getItem("token")}`;
  if (data.sseEvent){
    data.sseEvent.close();
  }
  data.sseEvent = new EventSource(url);
  data.sseEvent.onopen = () =>{
    console.log("handshake success");
  }
  data.sseEvent.onerror = (err)=>{
    console.log(err);
  }
  data.sseEvent.addEventListener("event_log",async function(res){
    let body = await JSON.parse(res.data);
    data.eventLogs = body.data.data;
    data.page =  body.data.cursor;
    data.total = Math.ceil(body.data.total / data.pageSize);
  })
}

onMounted(async()=>{

  data.form = {
    id:sessionStorage.getItem("id")??"",
    status:sessionStorage.getItem("status")??""
  };
  data.page = sessionStorage.getItem("page")??1;

  initEventSource();

})

onUnmounted(()=>{
  data.sseEvent.close();
})

const {eventLogs,form,page,total,cursor} = toRefs(data);

</script>
<style scoped>
.event{
  transition: opacity 0.5s ease;
  opacity: 1;
}
</style>