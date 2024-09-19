<template>
  <div class="event">
    <div class="container-fluid">

        <div class="mb-3 row">
          <div class="col-2">
            <div class="row">
            <label class="col-sm-2 col-form-label">Id:</label>
            <div class="col-sm-10">
              <input type="text" class="form-control" id="staticEmail" value="">
            </div>
            </div>
          </div>

          <div class="col-2">
            <div class="row">
            <label class="col-sm-2 col-form-label">Status:</label>
            <div class="col-sm-10">
              <select class="form-select" aria-label="Default select" style="cursor: pointer">
                <option selected>Open this select</option>
                <option value="1">Published</option>
                <option value="2">Success</option>
                <option value="3">Failed</option>
              </select>
            </div>
            </div>
          </div>
          <div class="col-2">
            <div class="col-auto">
              <button type="submit" class="btn btn-primary mb-3">Search</button>
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
            <th scope="col">Payload</th>
            <th scope="col">Action</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <th scope="row">1</th>
            <td>Mark</td>
            <td>Otto</td>
            <td>@mdo</td>
            <td>@mdo</td>
            <td>@mdo</td>
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
          <tr>
            <th scope="row">2</th>
            <td>Jacob</td>
            <td>Thornton</td>
            <td>@fat</td>
            <td>@mdo</td>
            <td>@mdo</td>
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
<style scoped>
.event{
  transition: opacity 0.5s ease;
  opacity: 1;
}
</style>