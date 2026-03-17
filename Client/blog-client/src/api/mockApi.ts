const blogs = [
  {
    blogID: 77,
    author: {
      authorID: "00000000-0000-0000-0000-000000000002",
      nickname: "admin",
      fullName: "adminadmin",
      email: "",
    },
    urlSlug: "new-blog",
    title: "New blog",
    content: "ABC",
    active: "true",
    createdAt: "2026-03-15T23:00:52.522928+09:00",
    createdBy: "00000000-0000-0000-0000-000000000002",
    updatedAt: "2026-03-15T23:00:52.522928+09:00",
    updatedBy: "00000000-0000-0000-0000-000000000002",
    deletedAt: null,
  },
  {
    blogID: 78,
    author: {
      authorID: "00000000-0000-0000-0000-000000000002",
      nickname: "admin",
      fullName: "adminadmin",
      email: "",
    },
    urlSlug: "my-new-blog",
    title: "My new Blog !@#",
    content: "asdasd",
    active: "true",
    createdAt: "2026-03-15T23:31:38.786923+09:00",
    createdBy: "00000000-0000-0000-0000-000000000002",
    updatedAt: "2026-03-15T23:31:38.786923+09:00",
    updatedBy: "00000000-0000-0000-0000-000000000002",
    deletedAt: null,
  },
  {
    blogID: 82,
    author: {
      authorID: "00000000-0000-0000-0000-000000000002",
      nickname: "admin",
      fullName: "adminadmin",
      email: "",
    },
    urlSlug: "test-blog",
    title: "test blog",
    content: "asdasdsdad",
    active: "true",
    createdAt: "2026-03-16T22:55:00.066542+09:00",
    createdBy: "00000000-0000-0000-0000-000000000002",
    updatedAt: "2026-03-16T22:55:00.066542+09:00",
    updatedBy: "00000000-0000-0000-0000-000000000002",
    deletedAt: null,
  },
];

const loginResponse = {
  userID: "00000000-0000-0000-0000-000000000002",
  firstName: "admin",
  lastName: "admin",
  roles: ["admin"],
  email: "",
  active: "normal",
  profileData: {},
  accessToken:
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDIiLCJVc2VybmFtZSI6InJvb3QiLCJGaXJzdE5hbWUiOiJhZG1pbiIsIkxhc3ROYW1lIjoiYWRtaW4iLCJSb2xlcyI6WyJhZG1pbiJdLCJUb2tlblZlcnNpb24iOjAsImlzcyI6IkJsb2dTZXJ2ZXIiLCJleHAiOjE3NzM3MTk3MDIsImlhdCI6MTc3MzcxNzkwMn0.rZmAqHyZeXnq64ImcPrcZenXYeww0claW9nIgDphj00",
};

const blogPOST = "Okay";

const dashboardData = { message: "Welcome to dashboard" };

const me = {
  userID: "00000000-0000-0000-0000-000000000002",
  firstName: "admin",
  lastName: "admin",
  roles: ["admin"],
  email: "",
  active: "normal",
  profileData: {},
  accessToken: "",
};

const updateUserData = { message: "OK" };

const notifications = [
  {
    notificationId: 49,
    userId: "",
    content: "A blog with title New blog has just been created",
    isRead: false,
    createdAt: "2026-03-15T22:57:16.271971+09:00",
    updatedAt: "0001-01-01T00:00:00Z",
    deletedAt: null,
  },
  {
    notificationId: 50,
    userId: "",
    content: "A blog with title New blog has just been created",
    isRead: false,
    createdAt: "2026-03-15T22:59:45.272324+09:00",
    updatedAt: "0001-01-01T00:00:00Z",
    deletedAt: null,
  },
  {
    notificationId: 51,
    userId: "",
    content: "A blog with title New blog has just been created",
    isRead: false,
    createdAt: "2026-03-15T23:00:53.162293+09:00",
    updatedAt: "0001-01-01T00:00:00Z",
    deletedAt: null,
  },
  {
    notificationId: 52,
    userId: "",
    content: "A blog with title My new Blog !@# has just been created",
    isRead: false,
    createdAt: "2026-03-15T23:31:39.011637+09:00",
    updatedAt: "0001-01-01T00:00:00Z",
    deletedAt: null,
  },
  {
    notificationId: 53,
    userId: "",
    content: "A blog with title test blog has just been created",
    isRead: false,
    createdAt: "2026-03-16T22:55:00.15002+09:00",
    updatedAt: "0001-01-01T00:00:00Z",
    deletedAt: null,
  },
  {
    notificationId: 54,
    userId: "",
    content: "A blog with title My first blog has just been created",
    isRead: false,
    createdAt: "2026-03-17T12:26:25.142313+09:00",
    updatedAt: "0001-01-01T00:00:00Z",
    deletedAt: null,
  },
];

const emailCode = { code: "123456" };

export {
  blogs,
  loginResponse,
  dashboardData,
  blogPOST,
  me,
  updateUserData,
  notifications,
  emailCode,
};
