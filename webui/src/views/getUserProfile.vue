<script>

export default {
  data: function () {
    return {
      userId: localStorage.getItem("Identifier"),
      profileId: localStorage.getItem("ProfileId"),
      userProfile: null,
      errorMessage: '',
      searchUsername: '',
      autocompleteResults: [],
      searchResultsData: [],
      isFollowing: false,
      isCurrentUser: false,
      IsBanned: false,
      selectedPhoto: null,
      loadingComments: false,
      newCommentText: '',
      newUsername: '',
    }
  },

  watch: {
    '$route': 'fetchUserProfileOnRouteChange',
  },

  async mounted() {
    await this.getUserProfile();
    document.addEventListener('click', this.closeAutocompleteList);
  },

  beforeUnmount() {
    // Rimuove il gestore di eventi quando il componente viene distrutto per evitare memory leaks
    document.removeEventListener('click', this.closeAutocompleteList);
  },
  methods: {

    async searchUser() {
      try {
        const response = await this.$axios.get(`${this.$url}/user/profile/${this.userId}/search/${this.searchUsername}`, {
          headers: {
            Authorization: 'Bearer ' + this.userId,
          },
        });

        if (response.status !== 200) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        // Pass the search results to the SearchResults component
        this.searchResultsData = response.data;
        this.errorMessage = ''; 
        localStorage.setItem("searchResultsData",JSON.stringify(this.searchResultsData));
        this.$router.push({ name: 'search-results'});

      } catch (error) {
        console.log('Error fetching search results:', error);
        this.errorMessage = 'Utente non trovato';
        this.userProfile = null; 
      }
    },
    async searchUserAutocomplete() {
    try {
      if (this.searchUsername.trim() === '') {
        this.autocompleteResults = [];
        return;
      }

      const response = await this.$axios.get(`${this.$url}/user/profile/${this.userId}/search/${this.searchUsername}`, {
        headers: {
          Authorization: 'Bearer ' + this.userId
        }
      });

      console.log('Searching user for autocomplete...');

      if (response.status !== 200) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      // Aggiorna la lista degli utenti per l'autocompletamento
      this.autocompleteResults = response.data;
    } catch (error) {
      console.error('Error fetching user autocomplete:', error);
    }
  },

    async goToUserProfile(selectedUserId) {
      const userId = localStorage.getItem("Identifier");
      localStorage.setItem("ProfileId", selectedUserId);
      this.$router.push({ path: "/users/"+ userId +"/profile/" + selectedUserId });
      this.autocompleteResults = [];
    },

    closeAutocompleteList(event) {
      // Chiude la lista degli utenti se il clic non è avvenuto sulla barra di ricerca o sulla lista stessa
      const searchInput = this.$refs.searchInput;
      const autocompleteList = this.$refs.autocompleteList;

      if (searchInput && autocompleteList && !searchInput.contains(event.target) && !autocompleteList.contains(event.target)) {
        this.autocompleteResults = [];
      }
    },

    async setMyUserName() {
      try {
        const response = await this.$axios.put(`${this.$url}/user/${this.userId}`, JSON.stringify(this.newUsername), {
          headers: { Authorization: 'Bearer ' + this.userId },
        });

        if (response.status === 201) {
          await this.getUserProfile();
        } else {
          console.error('Error changing username:', response.status, response.data);
        }

      } catch (error) {
        console.error('Error changing username:', error);
      }
    },

    async followUser() {
      try {
        const followResponse = await this.$axios.put(`${this.$url}/user/${this.userId}/follow/${this.profileId}`, null,{
          headers: {
            Authorization: 'Bearer ' + this.userId
          }
        });

        if (followResponse.status === 201) {
          this.isFollowing = true;
          await this.getUserProfile();
        }
      } catch (error) {
        console.error('Error following user:', error);
      }
    },

    async unfollowUser() {
      try {
        const unfollowResponse = await this.$axios.delete(`${this.$url}/user/${this.userId}/follow/${this.profileId}`, {
          headers: {
            Authorization: 'Bearer ' + this.userId
          }
        });

        if (unfollowResponse.status === 204) {
          this.isFollowing = false;
          await this.getUserProfile();
        }
      } catch (error) {
        console.error('Error unfollowing user:', error);
      }
    },

    async banUser() {
     try {
        // Chiamata API per bannare l'utente
       const banresponse = await this.$axios.post(`${this.$url}/user/${this.userId}/ban/${this.profileId}`, null, {
          headers: {
            Authorization: 'Bearer ' + this.userId,
          },
        });

        if(banresponse.status === 201) {
          this.IsBanned = true;
          this.isFollowing = false; 
          await this.getUserProfile();
        }
      } catch (error) {
        console.error('Error banning user:', error);
      }
    },

    async unbanUser() {
      try {
        // Chiamata API per sbannare l'utente
        const unbanresp = await this.$axios.delete(`${this.$url}/user/${this.userId}/ban/${this.profileId}`, {
          headers: {
            Authorization: 'Bearer ' + this.userId,
          },
        });

        if(unbanresp.status === 204) {
          this.IsBanned = false;
          await this.getUserProfile();
        }
      } catch (error) {
        console.error('Error unbanning user:', error);
      }
    },

    async getUserProfile() {
      try {
      const response = await this.$axios.get(`${this.$url}/users/${this.userId}/profile/${this.profileId}`, {
        headers: {
          Authorization: 'Bearer ' + this.userId
        }
      });

      if (response.status != 201 && response.status != 200) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      this.userProfile = response.data;
      console.log(this.userProfile);

      if (this.userId === this.profileId) {
        this.isCurrentUser = true;
      } else {
        const bancheck = await this.checkBanned();
        console.log(bancheck);
        if (!bancheck) {
            this.IsBanned = false;
            const followStatusResponse = await this.checkFollower();
            if (followStatusResponse) {
              this.isFollowing = true;
            }
          } else {
            this.IsBanned = true;
          }
        }
      } catch (error) {
        console.error('Error fetching user profile:', error);
        this.errorMessage = 'Unable to fetch user profile';
      }
    },

    async checkBanned(){
      try {
        const response = await this.$axios.get(`${this.$url}/users/${this.userId}/bans/${this.profileId}`, {
          headers: { Authorization: 'Bearer ' + this.userId },
        });
        console.log(response.data);
        if(response.data === true){
          return true;
        }
        return false;
      }
      catch (error) {
        console.error('Error checking ban:', error);
      }
    },

    async checkFollower(){
      try {
        const response = await this.$axios.get(`${this.$url}/users/${this.userId}/follows/${this.profileId}`, {
          headers: { Authorization: 'Bearer ' + this.userId },
        });
        console.log(response.data);
        if(response.data === true){
          return true;
        }
        return false;
      }
      catch (error) {
        console.error('Error checking follower:', error);
      }

    },

    async fetchUserProfileOnRouteChange() {
      // Recupera i nuovi parametri dell'URL
      const newUserId = localStorage.getItem("Identifier");
      const newProfileId = localStorage.getItem("ProfileId");

      // Aggiorna i dati del profilo solo se i parametri sono diversi
      if (newUserId !== this.userId || newProfileId !== this.profileId) {
        this.userId = newUserId;
        this.profileId = newProfileId;
        this.isCurrentUser = newUserId === newProfileId;
        await this.getUserProfile();
      }
    },

    openFileInput() {
      // Apre l'input di tipo file quando l'utente clicca sul pulsante "Carica Foto"
      this.$refs.fileInput.click();
    },

    handleFileChange(event) {
      // Gestisce l'evento di cambio dell'input di tipo file
      const fileInput = event.target;
      const selectedFile = fileInput.files[0];

      if (selectedFile) {
        this.uploadPhoto(selectedFile);
      }
    },

  async uploadPhoto(file) {
      try {
        // Usa FormData per inviare il file al backend
        console.log(file)
        const formData = new FormData();
        formData.append('file', file);
        console.log(formData)

        // Effettua la chiamata al backend per caricare la foto
        const response = await this.$axios.post(`${this.$url}/user/${this.userId}`, formData, {
          headers: {
            Authorization: 'Bearer ' + this.userId,
            'Content-Type': 'multipart/form-data',
          },
          
        });

        // Verifica lo stato della risposta
        if (response.status === 201) {
          // Ricarica il profilo utente
          await this.getUserProfile();
        } else {
          console.error('Error uploading photo:', response.status, response.data);
        }
      } catch (error) {
        console.error('Error uploading photo:', error);
      }
    },

      getPhotoUrl(imageData) {

      var buffer = Uint8Array.from(atob(imageData), c => c.charCodeAt(0)).buffer;

      // Creare un oggetto Blob dall'array buffer
      var blob = new Blob([buffer], { type: 'image/png' });

      // Creare un URL blob dal Blob creato
      var blobUrl = URL.createObjectURL(blob);

      // Assegnare l'URL blob all'attributo src dell'immagine
      return blobUrl;
      
    },

      
  async getComments(photo) {
      if (this.selectedPhoto === photo) {
        this.selectedPhoto = null;
      } else {
        this.selectedPhoto = photo;
        this.loadingComments = true;
        try {
          const commentsResponse = await this.$axios.get(`${this.$url}/photos/${photo.Id}/comments`, {
            headers: {
              Authorization: 'Bearer ' + this.userId,
            },
          });

          if (commentsResponse.status === 200) {
            this.selectedPhoto.Comments = commentsResponse.data;
          } else {
            console.error('Error fetching comments:', commentsResponse.status, commentsResponse.data);
          }
        } catch (error) {
          console.error('Error fetching comments:', error);
        }
        this.loadingComments = false;
      }
    },

    async deletePhoto(photo) {
      try {
        const response = await this.$axios.delete(`${this.$url}/user/${this.userId}/photos/${photo.Id}`, {
          headers: {
            Authorization: 'Bearer ' + this.userId,
          },
        });

        if (response.status === 204) {
          await this.getUserProfile();
        } else {
          console.error('Error deleting photo:', response.status, response.data);
        }
      } catch (error) {
        console.error('Error deleting photo:', error);
      }
    },

    async commentPhoto(photo) {
      try {
        // Effettua una chiamata API per aggiungere un nuovo commento
        const response = await this.$axios.post(
          `${this.$url}/user/${this.userId}/photos/${photo.Id}`, JSON.stringify(this.newCommentText),
          {
            headers: {
              Authorization: 'Bearer ' + this.userId,
            },
          }
        );

        if (response.status === 201) {
          await this.getUserProfile();
          // await this.getComments(photo);
          this.newCommentText = '';
        } else {
          console.error('Error adding new comment:', response.status, response.data);
        }
      } catch (error) {
        console.error('Error adding new comment:', error);
      }
    },

    async uncommentPhoto(comment) {
      try {
        const response = await this.$axios.delete(`${this.$url}/user/${this.userId}/photos/${comment.PhotoId}/comment/${comment.Id}`, {
          headers: {
            Authorization: 'Bearer ' + this.userId,
          },
        });

        if (response.status === 204) {
          await this.getUserProfile();
        } else {
          console.error('Error deleting comment:', response.status, response.data);
        }
      } catch (error) {
        console.error('Error deleting comment:', error);
      }
    },

    async likePhoto(photo) {
      if(photo.UserID == this.userId){
        return;
      }
      try {
        const checkLikeResponse = await this.checkLiked(photo);
        console.log(checkLikeResponse);
        if (checkLikeResponse == true){
          await this.unlikePhoto(photo);
        } else {
          const response = await this.$axios.post(`${this.$url}/user/${this.userId}/likes/${photo.Id}`, null, {
            headers: {
              Authorization: 'Bearer ' + this.userId
            },
          });

          if (response.status === 201) {
            await this.getUserProfile();
          } else {
            console.error('Error liking photo:', response.status, response.data);
          }
        }
      } catch (error) {
        console.error('Error liking photo:', error);
      }
    },

    async checkLiked(photo){
      try {
        const response = await this.$axios.get(`${this.$url}/users/${this.userId}/likes/${photo.Id}`, {
          headers: { Authorization: 'Bearer ' + this.userId },
        });
        console.log(response.data);
        if(response.data === true){
          return true;
        }
        return false;
      }
      catch (error) {
        console.error('Error checking like:', error);
      }
    },


    async unlikePhoto(photo){
      try {
        const response = await this.$axios.delete(`${this.$url}/user/${this.userId}/likes/${photo.Id}`, {
            headers: {
              Authorization: 'Bearer ' + this.userId,
            },
          });

        if (response.status === 204) {
          await this.getUserProfile();
        } else {
          console.error('Error deleting like:', response.status, response.data);
        }
      } catch (error) {
        console.error('Error deleting like:', error);
      }
    },
  },
}

</script>

<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col">
        <h1 class="text-center mb-4">User Profile - {{ userProfile ? userProfile.Username : '' }}</h1>
        <div class="input-group mb-3">
          <input ref="searchInput" type="text" class="form-control rounded-pill" v-model="searchUsername" placeholder="Cerca utente" @input="searchUserAutocomplete()"/>
          <div class="input-group-append">
            <button class="btn btn-outline-secondary btn-lg rounded-pill btn-hover" @click="searchUser" type="button">Cerca</button>
          </div>
        </div>
        <div v-if="autocompleteResults.length > 0" class="autocomplete-results" ref="autocompleteList"> <!-- Aggiunge un riferimento alla lista degli utenti -->
          <ul class="list-group">
            <li v-for="result in autocompleteResults" :key="result.Id" @click="goToUserProfile(result.Id)" class="list-group-item">{{ result.Username }}</li>
          </ul>
        </div>
        <div v-if="errorMessage" class="alert alert-danger mt-2" role="alert">
          <strong>{{ errorMessage }}</strong>
        </div>
        <div v-if="userProfile" class="text-center bg-dark p-4 rounded">
          <div class="d-flex justify-content-around mb-4">
            <p class="text-light h4">{{ userProfile.Username }}</p>
            <p class="text-light h5">Post: {{ userProfile.Photo_count }}</p>
            <p class="text-light h5">Follower: {{ userProfile.Follower_count }}</p>
            <p class="text-light h5">Seguiti: {{ userProfile.Following_count }}</p>
            <button v-if="!isCurrentUser && !isFollowing && !IsBanned" @click="followUser" class="btn btn-primary btn-follow">Follow</button>
            <button v-if="!isCurrentUser && isFollowing" @click="unfollowUser" class="btn btn-danger">Unfollow</button>
            <button v-if="!isCurrentUser && IsBanned" @click="unbanUser" class="btn btn-success">Unban</button>
            <button v-if="!isCurrentUser && !IsBanned" @click="banUser" class="btn btn-danger">Ban</button>
            <input type="file" ref="fileInput" @change="handleFileChange" style="display: none">
            <button v-if="isCurrentUser" @click="openFileInput" class="btn btn-success">Carica Foto</button>
            <div class="new-comment-container" v-if="isCurrentUser">
              <input v-model="newUsername" type="text" placeholder="Nuovo Username">
              <button @click="setMyUserName" class="btn btn-primary btn-follow">Cambia</button>
            </div>
          </div>
          <div class="scrollable-list">
            <div v-for="photo in userProfile.Photos" :key="photo.Id" class="mb-4" :class="{ 'post-bigger': showComments }">
              <div class="post-header">
                  <div style="flex-grow: 1;">
                    <span class="author-name">{{ userProfile.Username }}</span>
                    <p class="text-light small">Uploaded on: {{ photo.DateUploaded }}</p>
                  </div>
                  <button v-if="isCurrentUser" @click="deletePhoto(photo)" class="btn btn-danger btn-sm delete-post">Delete</button>
              </div>
              <img :src="getPhotoUrl(photo.ImageData)" class="img-fluid rounded" alt="Responsive image">
              <div class="post-details">
              <div class="like-section" @click="likePhoto(photo)">
                <span class="like-icon">❤️</span>
                <span class="like-counter">{{ photo.NumLikes }}</span>
              </div>
                <div class="comments-section">
                  <span class="comment-icon" @click="getComments(photo)">💬</span>
                  <span class="comment-counter">{{ photo.NumComments }}</span>
                  <div v-if="!loadingComments && selectedPhoto && selectedPhoto === photo" class="comment-list-container">
                    <div v-for="comment in selectedPhoto.Comments" :key="comment.Id" class="comment">
                      <span class="comment-author">{{ comment.Username }}</span>
                      <span class="comment-text">{{ comment.Comment }}</span>
                      <button v-if="comment.UserID == this.userId" @click="uncommentPhoto(comment)" class="btn btn-danger btn-sm btn-delete" title="Elimina il tuo commento">X</button>
                    </div>
                    <div class="new-comment-container">
                      <textarea v-model="newCommentText" placeholder="Inserisci un nuovo commento"></textarea>
                      <button @click="commentPhoto(photo)">Invia</button>
                    </div>
                  </div>
                </div>
              </div>  
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.btn-hover:hover {
  background-color: #007bff; /* Colore blu di Bootstrap */
  color: #fff; /* Colore del testo quando il bottone è in hover */
}

.autocomplete-results {
    position: absolute;
    width: 70%;
    background-color: #fff;
    border: 1px solid #ddd;
    max-height: 150px;
    overflow-y: auto; /* Abilita lo scrolling verticale se necessario */
    z-index: 1000;
    margin-top: 2px;
}

.autocomplete-results li {
  cursor: pointer;
}

.btn-follow {
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px;
}

.center-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.scrollable-list {
  width: 50%;
  overflow-y: scroll;
  display: flex;
  flex-direction: column;
}
  .post {
  position: relative;
  width: 100%;
  max-width: 400px; /* Adjust the max-width as needed */
  background-color: #f0f0f0;
  margin: 15px auto; /* Adjust the margin to center posts and provide spacing */
  padding: 15px; /* Add padding to the posts */
  display: flex;
  flex-direction: column;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s;
  border-radius: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.post:hover {
  transform: scale(1.02);
}

.post-header {
  background-color: #333;
  color: #fff;
  padding: 10px 10px 0 10px;
  text-align: left;
  border-top-left-radius: 10px;
  border-top-right-radius: 10px;
  display: flex;
  flex-wrap: wrap;
}

.author-name {
  font-weight: bold;
  flex-grow: 1;
  align-self: flex-start;
}

.delete-post{
  flex-grow: 0;
  flex-shrink: 0;
  height: 80%;
}

.image {
  width: 100%;
  object-fit: cover;
  border-bottom-left-radius: 10px;
  border-bottom-right-radius: 10px;
}

.post-details {
  display: flex;
  flex-direction: row;
  align-items: center;
  background-color: #333;
  color: #fff;
  transition: transform 0.3s, max-height 0.3s;
  padding: 10px; /* Add padding to the post details */
}

  .like-section,
  .comments-section {
    display: flex;
    align-items: center;
    padding: 10px;
    cursor: pointer;
  }

  .like-icon,
  .comment-icon {
    margin-right: 5px;
    font-size: 18px;
  }

.like-counter,
.comment-counter {
  margin-right: 5px;
  font-size: 16px;
}

.comment-list-container {
  width: 300px;
  padding: 15px;
  background-color: #333;
  border: 1px solid #ddd;
  position: relative; /* Change from absolute to relative */
  margin-top: 10px; /* Add margin to position it below the post */
  overflow-y: auto;
  border-radius: 8px;
}

.comment {
  margin-bottom: 8px;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  display: flex;
  flex-direction: row;
  align-items: flex-start;
}

.btn-delete{
  flex-grow: 0;
  flex-shrink: 0;
}

.comment-author {
  flex-grow: 1;
  font-weight: bold;
  margin-right: 5px;
  align-self: flex-start;
}

.comment-text {
  flex-grow: 6;
  align-self: flex-start;
  text-align: left;
}

.post-bigger {
  width: 600px;
}

.new-comment-container {
  margin-top: 10px;
}

.new-comment-container textarea {
  width: 100%;
  resize: vertical; /* Consente la regolazione verticale della textarea */
  margin-bottom: 5px;
}

.new-comment-container button {
  background-color: #007bff;
  color: #fff;
  border: none;
  padding: 5px 10px;
  cursor: pointer;
  border-radius: 5px;
}



</style>