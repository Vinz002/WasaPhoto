<script>
export default {
  data: function () {
    return {
      searchResultsData: [],
    };
  },
  mounted() {
  console.log('SearchResults component mounted');
  const searchData = localStorage.getItem("searchResultsData");

  if (searchData) {
    this.searchResultsData = JSON.parse(searchData);
  }
  console.log(this.searchResultsData);
},
  
  methods: {
    goToUserProfile(selectedUserId) {
      const userId = localStorage.getItem("Identifier");
      localStorage.setItem("ProfileId", selectedUserId);
      this.$router.push({ path: "/users/"+ userId +"/profile/" + selectedUserId });
    },
  },
};
</script>

<template>
  <div class="container mt-4">
    <h1 class="text-center mb-4">Search Results</h1>
    <div v-if="searchResultsData.length > 0" class="autocomplete-results" ref="autocompleteList">
          <ul class="list-group">
              <li v-for="result in this.searchResultsData" :key="result.Id" @click="goToUserProfile(result.Id)" class="list-group-item">{{ result.Username }}</li>
          </ul>
     </div>
  </div>
</template>

<style scoped>
  .autocomplete-results {
      border-radius: 5px;
      cursor: pointer;
  }

</style>