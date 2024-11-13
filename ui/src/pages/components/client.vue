<template>
  <div class="card">
    <div class="card-header">
      Client
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="connected" class="col-sm-7 col-form-label">Connected Clients</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="connected" :value="clients.connected_clients">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="blocked" class="col-sm-7 col-form-label">Blocked Clients</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="blocked" :value="clients.blocked_clients">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="maxInBuffer" class="col-sm-7 col-form-label">Client Recent Max Input Buffer</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="maxInBuffer" :value="clients.client_recent_max_input_buffer">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="maxOutBuffer" class="col-sm-7 col-form-label">Client Recent Max Output Buffer</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="maxOutBuffer" :value="clients.client_recent_max_output_buffer">
      </div>
    </div>
  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  clients:{}
})

</script>