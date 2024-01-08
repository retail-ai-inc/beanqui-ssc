<template>
    <div>
      <Pagination :page="page" :total="total" @changePage="changePage"/>

<!--      <div class="accordion" id="ui-accordion">-->
<!--        <div class="accordion-item" v-for="(item, key) in queues" :key="key">-->
<!--          <h2 class="accordion-header">-->
<!--            <button style="font-weight: bold" class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="setId(key)" aria-expanded="true" :aria-controls="key">-->
<!--              {{key}}-->
<!--            </button>-->
<!--          </h2>-->
<!--          <div :id="key" class="accordion-collapse collapse show" data-bs-parent="#ui-accordion">-->
<!--            <div class="accordion-body">-->
<!--              <table class="table table-striped">-->
<!--                <thead>-->
<!--                <tr>-->
<!--                  <th scope="col">Group</th>-->
<!--                  <th scope="col">Queue</th>-->
<!--                  <th scope="col">State</th>-->
<!--                  <th scope="col">Size</th>-->
<!--                  <th scope="col">Memory usage</th>-->
<!--                  <th scope="col">Processed</th>-->
<!--                  &lt;!&ndash;                    <th scope="col">Failed</th>&ndash;&gt;-->
<!--                  &lt;!&ndash;                    <th scope="col">Error rate</th>&ndash;&gt;-->
<!--                  <th scope="col">Action</th>-->
<!--                </tr>-->
<!--                </thead>-->
<!--                <tbody>-->
<!--                <tr v-for="(d, k) in item" :key="k">-->

<!--                  <th scope="row">{{d.group}}</th>-->

<!--                  <th scope="row">{{ d.queue }}</th>-->
<!--                  <td :class="d.state == 'Run' ? 'text-success-emphasis' : 'text-danger-emphasis'">{{ d.state }}</td>-->
<!--                  <td>{{ d.size }}</td>-->
<!--                  <td>{{ d.memory }}</td>-->
<!--                  <td>{{ d.process }}</td>-->
<!--                  &lt;!&ndash;                    <td>{{ item.fail }}</td>&ndash;&gt;-->
<!--                  &lt;!&ndash;                    <td>{{ item.errRate }}</td>&ndash;&gt;-->
<!--                  <td>-->
<!--                    <div class="btn-group" role="group" aria-label="Button group with nested dropdown">-->
<!--                      <div class="btn-group" role="group">-->
<!--                        <button type="button" class="btn btn-primary dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">-->
<!--                          Actions-->
<!--                        </button>-->
<!--                        <ul class="dropdown-menu">-->
<!--                          <li><a class="dropdown-item" href="#">Delete</a></li>-->
<!--                          <li><a class="dropdown-item" href="#">Pause</a></li>-->
<!--                        </ul>-->
<!--                      </div>-->
<!--                    </div>-->
<!--                  </td>-->
<!--                </tr>-->
<!--                </tbody>-->
<!--              </table>-->
<!--            </div>-->
<!--          </div>-->
<!--        </div>-->
<!--      </div>-->


        <table class="table table-striped">
            <thead>
                <tr>
                  <th scope="col">Group</th>
                    <th scope="col">Queue</th>
                    <th scope="col">State</th>
                    <th scope="col">Size</th>
                    <th scope="col">Memory usage</th>
                    <th scope="col">Processed</th>
                    <th scope="col">Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(item, key) in queues" :key="key">

                    <th scope="row">{{key}}</th>
                    <th scope="row">{{ item.queue }}</th>
                    <td :class="item.state == 'Run' ? 'text-success-emphasis' : 'text-danger-emphasis'">{{ item.state }}</td>
                    <td>{{ item.size }}</td>
                    <td>{{ item.memory }}</td>
                    <td>{{ item.process }}</td>
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
      <Pagination :page="page" :total="total" @changePage="changePage"/>
    </div>
</template>
  
  
<script setup>

import { reactive,onMounted,toRefs,onUnmounted } from "vue";
import request  from "request";
import Pagination from "./components/pagination.vue";

let pageSize = 10;
let data = reactive({
  queues:[],
  page:1,
  total:1
})

function getQueue(page,pageSize){
  return request.get("queue",{"params":{"page":page,"pageSize":pageSize}});
}
onMounted(async ()=>{
  let queue = await getQueue(data.page,10);
  data.queues = {...queue.data};
})
async function changePage(page){
  let queue = await getQueue(page,10);
  data.queues = {...queue.data.data};
  data.total = Math.ceil(queue.data.total / 10);
  data.page = page;
}
function setId(id){
  return "#"+id;
}
const {queues,page,total} = toRefs(data);
</script>
  
<style scoped>
.table .text-success-emphasis {
    color: var(--bs-green) !important;
}
.table-striped th{
  font-weight: 400 !important;
}
.table .text-danger-emphasis {
    color: var(--bs-danger) !important;
}
</style>
  
  