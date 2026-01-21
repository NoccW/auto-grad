import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Home from "../views/Home.vue";
import Upload from "../views/Upload.vue";
import Result from "../views/Result.vue";
import History from "../views/History.vue";
import TeacherHome from "../views/TeacherHome.vue";
import TaskConfig from "../views/TaskConfig.vue";
import TaskMonitor from "../views/TaskMonitor.vue";
import TaskAnalysis from "../views/TaskAnalysis.vue";
import TaskHistory from "../views/TaskHistory.vue";
import ParentDashboard from "../views/ParentDashboard.vue";
import UnifiedDashboard from "../views/UnifiedDashboard.vue";

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
  // 统一仪表板 - 根据用户角色显示不同内容
  {
    path: "/",
    name: "UnifiedDashboard",
    component: UnifiedDashboard,
    meta: { requiresAuth: true },
  },
  // 家长端仪表板
  {
    path: "/parent",
    name: "ParentDashboard",
    component: ParentDashboard,
    meta: { requiresAuth: true, role: "parent" },
  },
  // 教师端仪表板
  {
    path: "/teacher",
    name: "TeacherHome",
    component: TeacherHome,
    meta: { requiresAuth: true, role: "teacher" },
  },
  // 家长端功能
  {
    path: "/upload",
    name: "Upload",
    component: Upload,
    meta: { requiresAuth: true, roles: ["parent", "teacher"] },
  },
  {
    path: "/result/:id",
    name: "Result",
    component: Result,
    meta: { requiresAuth: true, roles: ["parent", "teacher"] },
    props: true,
  },
  {
    path: "/history",
    name: "History",
    component: History,
    meta: { requiresAuth: true, roles: ["parent", "teacher"] },
  },
  // 教师端专用功能
  {
    path: "/task-config",
    name: "TaskConfig",
    component: TaskConfig,
    meta: { requiresAuth: true, role: "teacher" },
  },
  {
    path: "/task-monitor/:id",
    name: "TaskMonitor",
    component: TaskMonitor,
    meta: { requiresAuth: true, role: "teacher" },
    props: true,
  },
  {
    path: "/task-analysis/:id",
    name: "TaskAnalysis",
    component: TaskAnalysis,
    meta: { requiresAuth: true, role: "teacher" },
    props: true,
  },
  {
    path: "/task-history",
    name: "TaskHistory",
    component: TaskHistory,
    meta: { requiresAuth: true, role: "teacher" },
  },
  // 兼容旧路由
  {
    path: "/teacher-home",
    redirect: "/teacher",
  },
  {
    path: "/parent-dashboard",
    redirect: "/parent",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫 - 检查用户权限
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem("token");
  const userRole = localStorage.getItem("userRole");

  // 检查是否需要认证
  if (to.meta.requiresAuth && !isAuthenticated) {
    next("/login");
    return;
  }

  // 检查角色权限
  if (to.meta.role && to.meta.role !== userRole) {
    // 如果用户角色不匹配，重定向到对应的仪表板
    if (userRole === "teacher") {
      next("/teacher");
    } else if (userRole === "parent") {
      next("/parent");
    } else {
      next("/login");
    }
    return;
  }

  // 检查多角色权限
  if (to.meta.roles && to.meta.roles.length > 0) {
    if (!to.meta.roles.includes(userRole)) {
      // 如果用户没有权限，重定向到对应的仪表板
      if (userRole === "teacher") {
        next("/teacher");
      } else if (userRole === "parent") {
        next("/parent");
      } else {
        next("/login");
      }
      return;
    }
  }

  next();
});

export default router;
