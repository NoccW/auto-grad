import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Home from "../views/Home.vue";
import Upload from "../views/Upload.vue";
import Result from "../views/Result.vue";
import History from "../views/History.vue";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: { requiresAuth: false },
  },
  {
    path: "/register",
    name: "Register",
    component: Register,
    meta: { requiresAuth: false },
  },
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: { requiresAuth: true },
  },
  {
    path: "/upload",
    name: "Upload",
    component: Upload,
    meta: { requiresAuth: true },
  },
  {
    path: "/result/:id",
    name: "Result",
    component: Result,
    meta: { requiresAuth: true },
    props: true,
  },
  {
    path: "/history",
    name: "History",
    component: History,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
