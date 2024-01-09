<template>
    <div>

      <div class="accordion" id="schedule-ui-accordion">
        <div class="accordion-item" v-for="(item, key) in schedule" :key="key" style="margin-bottom: 15px">
          <h2 class="accordion-header">
            <button style="font-weight: bold" class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="setScheduleId(key)" aria-expanded="true" :aria-controls="key">
              Channel:&nbsp;&nbsp;{{key}}
            </button>
          </h2>
          <div :id="key" class="accordion-collapse collapse show" data-bs-parent="#schedule-ui-accordion">
            <div class="accordion-body" style="padding: 0.5rem">
              <table class="table table-striped">
                <thead>
                <tr>
                  <th scope="col">Topic</th>
                  <th scope="col">State</th>
                  <th scope="col">Size</th>
                  <th scope="col">Memory usage</th>
                  <th scope="col">Processed</th>
                  <th scope="col">Action</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(d, k) in item" :key="k">
                  <th scope="row">{{ d.queue }}</th>
                  <td :class="d.state == 'Run' ? 'text-success-emphasis' : 'text-danger-emphasis'">{{ d.state }}</td>
                  <td>{{ d.size }}</td>
                  <td>{{ d.memory }}</td>
                  <td>{{ d.process }}</td>
                  <td>
                    <div class="btn-group" role="group" aria-label="Button group with nested dropdown">
                      <div class="btn-group" role="group">
                        <button type="button" class="btn btn-primary dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
                          Actions
                        </button>
                        <ul class="dropdown-menu">
                          <li><a class="dropdown-item" href="#">Delete</a></li>
                          <li><a class="dropdown-item" href="#">Pause</a></li>
                        </ul>
                      </div>
                    </div>
                  </td>
                </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
      <Pagination :page="page" :total="total" @changePage="changePage"/>
    </div>
</template>
  
  
<script setup>

import { reactive,toRefs,onMounted,onUnmounted } from "vue";
import request  from "request";
import Pagination from "./components/pagination.vue";

const data = reactive({
  page:1,
  total:1,
  schedule:[]
})
function getSchedule(page,pageSize){
  return request.get("schedule",{"params":{"page":page,"pageSize":pageSize}});
}
async function changePage(page){
  let schedule = await getSchedule(page,10);
  data.schedule = {...schedule.data};
  data.total = Math.ceil(schedule.data.total / 10);
  data.page = page;
}
onMounted(async ()=>{
  let schedule = await getSchedule(data.page,10);
  data.schedule = {...schedule.data};
  data.total = Math.ceil(schedule.data.total / 10);
})
function setScheduleId(id){
  return "#"+id;
}
const {page,total,schedule} = toRefs(data);
</script>
  
<style scoped>
.table .text-success-emphasis{
    color:var(--bs-green) !important;
}
.table .text-danger-emphasis{
    color:var(--bs-danger) !important;
}
</style>
  
  