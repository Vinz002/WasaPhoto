import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/doLogin.vue'
import stream from '../views/getMyStream.vue'
import Profile from  '../views/getUserProfile.vue'
import SearchResults from '../views/SearchResults.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path:'/user/stream/:userid', component: stream},
		{path: '/users/:userid/profile/:profileid', component: Profile},
		{path: '/search-results', name: 'search-results',component: SearchResults}
	]
})

export default router
