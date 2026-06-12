type PortfolioData = {
  name: string;
  title: string;
  about: string[];
  skills: {
    name: string;
    value: number;
  }[];
  projects: Project[];
  hobbyProjects: Project[];
  journey: {
    label: string;
  }[];
};

export interface Project {
  title: string;
  time: string;
  description: string;
  technologies: {
    server?: string[];
    frontend?: string[];
    backend?: string[];
    architecture?: string[];
    infrastructure?: string[];
  };
  roles: string[];
  details: string[];
}

export const portfolioData: PortfolioData = {
  name: "Cuong Tran",

  title: "Full Stack Engineer • Backend Specialist • DevOps Enthusiast",

  about: [
    "I'm a software engineer with over 2 years of experience in full-stack development, primarily working on migration projects that help modernize and maintain enterprise applications. My professional experience has mainly involved Java and Node.js development, along with frontend technologies such as React.",
    "Outside of work, I enjoy building personal projects to explore technologies and architectural approaches that I don't often get to use in my daily role. I have been developing backend services in Go and experimenting with modular, event-driven system design to better understand scalability and system boundaries. These projects also give me hands-on experience with Docker, AWS, GitHub Actions, and Nginx, allowing me to learn more about deployment, automation, and software delivery. I enjoy tackling technical challenges, evaluating trade-offs, and continuously expanding my knowledge across both application development and infrastructure-related topics.",
    "Passionate engineer focused on building scalable backend systems, cloud-native applications, and deployment automation.",
  ],

  skills: [
    { name: "Backend Development", value: 95 },
    { name: "DevOps", value: 40 },
    { name: "Cloud Infrastructure", value: 50 },
    { name: "Frontend Development", value: 80 },
  ],

  journey: [{ label: "Junior Developer" }, { label: "Backend Engineer" }],

  projects: [
    {
      title: "File And Document Management Systems",
      time: "2025/06/14",
      description: `
        Manage file with advance search patterns, and CRUD on documents
      `,
      technologies: { server: ["Java", "Primefaces", "Oracle DB"] },
      roles: ["coding", "testing"],
      details: [
        "Migrate a legacy projecct from an old language to Java with Primefaces framework for view and bean state management",
      ],
    },
    {
      title: "Membership Management Systems",
      time: "2024/06/14",
      description: `
        Manage membership with CRUD on coupon, profile, history, roles
      `,
      technologies: {
        server: ["ASP.NET 9", "C#", "MSTest"],
        frontend: ["Vue", "NuxtJs"],
      },
      roles: ["coding", "testing"],
      details: [
        "Migrate a legacy projecct from an old language to Java with Primefaces framework for view and bean state management",
      ],
    },
  ],

  hobbyProjects: [
    {
      title: "Content Management System",
      time: "2026/06/14",
      description: `Developed a production-ready blogging platform using a microservice-ready Modular Monolith architecture with Golang, Gin, PostgreSQL, and React. 
        Applied event-driven patterns including Saga orchestration, Outbox Pattern, and an Event Bus to ensure reliability and consistency while maintaining clear module boundaries. 
        Containerized the application with Docker and designed it for future cloud-native deployment, caching, and CI/CD automation.`,
      technologies: {
        backend: ["Golang", "Gin", "PostgreSQL", "sqlc", "REST API"],
        architecture: [
          "Modular Monolith",
          "Event-Driven Architecture",
          "Saga Pattern, Outbox Pattern",
          "Clean Architecture",
        ],
        frontend: ["React", "TanStack Query", "Axios", "Material UI"],
        infrastructure: [
          "Docker",
          "Nginx",
          "GitHub Actions (CI/CD)",
          "AWS (EC2, RDS, S3)",
          "Redis",
        ],
      },
      roles: ["design", "coding", "testing", "deploy"],
      details: [
        `Designed and developed a microservice-ready Modular Monolith Blog Platform using Golang, Gin, PostgreSQL, and React, emphasizing domain-driven modularity, scalability, and maintainability.`,
        `Implemented event-driven architecture with an in-memory event bus, Outbox Pattern, and Saga Pattern to ensure reliable inter-module communication and transactional consistency.`,
        `Built a high-performance backend using sqlc-generated type-safe SQL, avoiding ORM overhead while optimizing query performance and database control.`,
        `Containerized the application with Docker and designed for cloud-native deployment, observability, and future scalability through planned integration with AWS, Redis, Prometheus, and CI/CD pipelines.`,
      ],
    },
  ],
};
