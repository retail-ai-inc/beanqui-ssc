<template>
  <div class="card">
    <div class="card-header">
      Command
    </div>
    <table class="table">
      <thead>
      <tr>
        <th scope="col">Command</th>
        <th scope="col">Calls</th>
        <th scope="col">Usec</th>
        <th scope="col">UsecPerCall</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(v,k) in commands" :key="k">
        <th scope="row">{{v.command}}</th>
        <td>{{v.calls}}</td>
        <td>{{v.usec}}</td>
        <td>{{v.usec_per_call}}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  commands:Array
})

</script>