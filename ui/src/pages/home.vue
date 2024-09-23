<template>

  <div class="home">
    <div class="container-fluid text-center" style="padding: 0">
      <div class="row align-items-start" style="margin: 1.25rem 0;color:#fff;">
        <div class="col" style="background: #0d6efd;height:7.5rem;padding:1rem;">
          <div>Queue Total</div>
          <div style="font-weight: bold">
            <router-link to="/admin/queue" class="nav-link text-muted link-color" >{{queue_total}}</router-link></div>
        </div>
        <div class="col" style="background: #198754;height:7.5rem;padding:1rem;">
          <div>CPU Total</div>
          <div style="font-weight: bold">
            <router-link to="/admin/redis" class="nav-link text-muted link-color">{{num_cpu}}</router-link>
          </div>
        </div>
        <div class="col" style="background: #dc3545;height:7.5rem;padding:1rem;">
          <div>Fail Total</div>
          <div style="font-weight: bold">
            <router-link to="" class="nav-link text-muted link-color">{{fail_count}}</router-link>
          </div>
        </div>
        <div class="col" style="background: #20c997;height:7.5rem;padding:1rem;">
          <div>Success Total</div>
          <div style="font-weight: bold">
            <router-link to="" class="nav-link text-muted link-color">{{success_count}}</router-link>
          </div>
        </div>
        <div class="col" style="background: #343a40;height:7.5rem;padding:1rem;">
          <div>Total Payload</div>
          <div style="font-weight: bold">
            <router-link to="" class="nav-link text-muted link-color">{{db_size}}</router-link>
          </div>
        </div>
      </div>
    </div>
<!--    <div class="d-flex justify-content-between">-->
<!--      <div></div>-->
<!--      <div style="width:50%">-->
<!--        <v-chart class="chart" :option="barOption" />-->
<!--      </div>-->
<!--      <div style="width:50%">-->
<!--        <v-chart class="chart" :option="lineOption"/>-->
<!--      </div>-->
<!--    </div>-->

  </div>


</template>

<script setup>
import {ref,reactive,onMounted,toRefs,} from "vue";
import request  from "request";

let data = reactive({
  "queue_total":0,
  "db_size":0,
  "num_cpu":0,
  "fail_count":0,
  "success_count":0
})
function getTotal(){
  return request.get("dashboard");
}
onMounted(async ()=>{

  let total = await getTotal();
  Object.assign(data,total.data);

})

const barOption = ref({
  title: {
    text: 'Queue Size',
    left: 'left'
  },
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: [120, 200, 150, 80, 70, 110, 130],
      type: 'bar',
      showBackground: true,
      backgroundStyle: {
        color: 'rgba(180, 180, 180, 0.2)'
      }
    }
  ]
});
const lineOption = ref({
  title: {
    text: 'Tasks Processed'
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    bottom:"10",
    data: [
        {
          name:"succeed",
          lineStyle: {
            color: '#198754'
          },
          itemStyle:{
            color:"#198754"
          }
        },
        {
          name:"failed",
          lineStyle:{
            color:'#dc3545'
          },
          itemStyle: {
            color:'#dc3545'
          }
        }
    ]
  },
  toolbox: {
    feature: {
      //saveAsImage: {}
    }
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      name: 'succeed',
      type: 'line',
      lineStyle: {
        color: '#198754'
      },
      itemStyle:{
        color:"#198754"
      },
      data:[220, 182, 191, 234, 290, 330, 310]
    },
    {
      name: 'failed',
      type: 'line',
      lineStyle:{
        color:'#dc3545'
      },
      itemStyle:{
        color:"#dc3545"
      },
      data:[120, 132, 101, 134, 90, 230, 210]
    }
  ]
});
const {queue_total,db_size,num_cpu,fail_count,success_count} = toRefs(data);
</script>
<style scoped>
.home{
  transition: opacity 0.5s ease;
  opacity: 1;
}
.chart{
  width:100%;height:80vh;
}
.link-color{
  display: inline-block;
  color: #fff !important;
}
</style>