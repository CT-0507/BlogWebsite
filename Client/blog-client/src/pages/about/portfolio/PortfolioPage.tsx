import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Grid from "@mui/material/Grid";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import AnimatedSection from "@/components/AnimatedSection";
import SectionHeader from "@/components/SectionHeader";
import SkillProgress from "@/components/SkillProgress";
import JourneyStepper from "@/components/JourneyStepper";
import ProjectCard from "@/components/ProjectCard";

import { portfolioData } from "./data/portfolioData";
import { Link } from "react-router-dom";
import ContactForm from "./ContactForm";

export default function PortfolioPage() {
  const scrollTo = (id: string) => {
    const element = document.querySelector(id);
    element?.scrollIntoView({ behavior: "smooth" });
    console.log("Here");
  };
  return (
    <Container maxWidth="lg" sx={{ overflow: "hidden" }}>
      {/* HERO */}

      <Box minHeight="100vh" display="flex" alignItems="center">
        <Stack spacing={3}>
          <Typography variant="h2">{portfolioData.name}</Typography>

          <Typography variant="h5">{portfolioData.title}</Typography>

          <Stack direction="row" spacing={2}>
            <Button
              component={Link}
              to="#projects"
              variant="contained"
              onClick={() => scrollTo("#projects")}
            >
              View Projects
            </Button>

            <Button variant="outlined" onClick={() => scrollTo("#contact")}>
              Contact
            </Button>
          </Stack>
        </Stack>
      </Box>

      {/* ABOUT */}

      <Box sx={{ mt: 1 }}>
        <AnimatedSection direction="left">
          <Box>
            <SectionHeader title="About Me" subtitle="Professional Summary" />

            <Typography>{portfolioData.about}</Typography>
          </Box>
        </AnimatedSection>
      </Box>

      {/* TECHNICAL FOCUS */}

      <AnimatedSection direction="right">
        <Box>
          <SectionHeader title="Technical Focus" />

          <Grid container spacing={3}>
            {portfolioData.skills.map((skill) => (
              <Grid size={{ xs: 12, md: 6 }} key={skill.name}>
                <SkillProgress {...skill} />
              </Grid>
            ))}
          </Grid>
        </Box>
      </AnimatedSection>

      {/* JOURNEY */}

      <Box sx={{ mt: 1 }}>
        <AnimatedSection direction="left">
          <Box>
            <SectionHeader title="Career Journey" />

            <JourneyStepper items={portfolioData.journey} />
          </Box>
        </AnimatedSection>
      </Box>

      {/* PROJECTS */}

      <Box id="projects" sx={{ mt: 1 }}>
        <AnimatedSection direction="right">
          <Box>
            <SectionHeader title="Projects" />

            <Grid container spacing={3}>
              {portfolioData.projects.map((project) => (
                <Grid size={{ xs: 12, md: 6 }} key={project.title}>
                  <ProjectCard {...project} />
                </Grid>
              ))}
            </Grid>
          </Box>
        </AnimatedSection>
      </Box>

      {/* Hobby Project */}

      <Box sx={{ mt: 1 }}>
        <AnimatedSection direction="left">
          <Box>
            <SectionHeader title="Hobby Projects" />

            <Grid container spacing={3}>
              {portfolioData.hobbyProjects.map((project) => (
                <Grid size={{ xs: 12, md: 6 }} key={project.title}>
                  <ProjectCard {...project} />
                </Grid>
              ))}
            </Grid>
          </Box>
        </AnimatedSection>
      </Box>

      {/* CONTACT */}

      <ContactForm id="contact" />
    </Container>
  );
}
