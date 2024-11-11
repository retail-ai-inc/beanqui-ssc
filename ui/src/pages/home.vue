<template>

  <div class="home">
    <div class="container-fluid text-center">
      <Dashboard   :queue_total="queue_total"
                   :num_cpu="num_cpu"
                   :fail_count="fail_count"
                   :success_count="success_count"
                   :db_size="db_size"/>
    </div>

    <div class="container-fluid text-center">
      <div class="row justify-content-between">
        <div class="col-5">
          <Command :commands="commands" />
        </div>
        <div class="col-4">
          <Client :clients="clients" />
          <KeySpace :keyspace="keyspace" />
        </div>
        <div class="col-3">
          <Stats :stats="stats" />
        </div>
      </div>
    </div>


<!--    <div class="d-flex justify-content-between">-->

<!--        <div style="width:30%">-->
<!--          <v-chart class="chart" :option="barOption" />-->
<!--        </div>-->
<!--        <div style="width:30%">-->
<!--          <v-chart class="chart" :option="lineOption"/>-->
<!--        </div>-->
<!--        <div style="width: 30%">-->
<!--          <v-chart class="chart" :option="gaugeOption" />-->
<!--        </div>-->

<!--    </div>-->

  </div>


</template>

<script setup>
import {ref,reactive,onMounted,toRefs,} from "vue";
import Dashboard from "./components/dashboard.vue";
import Command from "./components/command.vue";
import Client from "./components/client.vue";
import KeySpace from "./components/keySpace.vue";
import Stats from "./components/stats.vue";


let data = reactive({
  "queue_total":0,
  "db_size":0,
  "num_cpu":0,
  "fail_count":0,
  "success_count":0,
  "commands":[],
  "clients":{},
  "stats":{},
  "keyspace":[]
})
function getTotal(){
  return request.get("dashboard");
}
onMounted(async ()=>{

  let total = await getTotal();

  Object.assign(data,total.data);
  data.commands = total.data.commands;
  data.clients = total.data.clients;
  data.stats = total.data.stats;
  data.keyspace = total.data.keyspace;
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

const gaugeOption = ref({
  title:{
    text:"Memory usage"
  },
  series: [
    {
      type: 'gauge',
      center: ['50%', '60%'],
      startAngle: 200,
      endAngle: -20,
      min: 0,
      max: 60,
      splitNumber: 5,
      itemStyle: {
        color: '#FFAB91'
      },
      progress: {
        show: true,
        width: 30
      },
      pointer: {
        show: false
      },
      axisLine: {
        lineStyle: {
          width: 30
        }
      },
      axisTick: {
        distance: -45,
        splitNumber: 5,
        lineStyle: {
          width: 2,
          color: '#999'
        }
      },
      splitLine: {
        distance: -52,
        length: 14,
        lineStyle: {
          width: 3,
          color: '#999'
        }
      },
      axisLabel: {
        distance: -20,
        color: '#999',
        fontSize: 20
      },
      anchor: {
        show: false
      },
      title: {
        show: false
      },
      detail: {
        valueAnimation: true,
        width: '60%',
        lineHeight: 40,
        borderRadius: 8,
        offsetCenter: [0, '-15%'],
        fontSize: 60,
        fontWeight: 'bolder',
        formatter: '{value} M',
        color: 'inherit'
      },
      data: [
        {
          value: 20
        }
      ]
    },
    {
      type: 'gauge',
      center: ['50%', '60%'],
      startAngle: 200,
      endAngle: -20,
      min: 0,
      max: 60,
      itemStyle: {
        color: '#FD7347'
      },
      progress: {
        show: true,
        width: 8
      },
      pointer: {
        show: false
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        show: false
      },
      axisLabel: {
        show: false
      },
      detail: {
        show: false
      },
      data: [
        {
          value: 20
        }
      ]
    }
  ]
});
const {queue_total,db_size,num_cpu,fail_count,success_count,commands,clients,stats,keyspace} = toRefs(data);
</script>
<style scoped>
.home{
  transition: opacity 0.5s ease;
  opacity: 1;
}
.chart{
  width:100%;height:80vh;
}

</style>