import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Chip from "@mui/material/Chip";
import Stack from "@mui/material/Stack";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import type { Project } from "@/pages/about/portfolio/data/portfolioData";
import Box from "@mui/material/Box";
import { capitalize } from "@mui/material";

export default function ProjectCard({
  title,
  description,
  technologies,
  roles,
  time,
}: Project) {
  return (
    <Card>
      <CardContent>
        <Stack spacing={2}>
          <Typography variant="h5">{title}</Typography>
          <Typography
            sx={{
              opacity: 0.7,
              fontSize: "0.9rem",
            }}
          >
            {time}
          </Typography>

          <Box>
            Roles:{" "}
            {roles.map((item, index) => (
              <Chip key={index} label={item} sx={{ mr: 1 }} />
            ))}
          </Box>

          <Stack spacing={1} flexWrap="wrap">
            {Object.entries(technologies).map((value, index) => (
              <Box key={index}>
                {`${capitalize(value[0])}: `}
                {value[1].map((item) => (
                  <Chip key={item} label={item} sx={{ mr: 1 }} />
                ))}
              </Box>
            ))}
          </Stack>

          <Accordion>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
              Description
            </AccordionSummary>

            <AccordionDetails>
              <Typography align="justify">{description}</Typography>
            </AccordionDetails>
          </Accordion>
        </Stack>
      </CardContent>
    </Card>
  );
}
