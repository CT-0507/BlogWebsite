import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import Container from "@mui/material/Container";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import { Link as RouterLink, useSearchParams } from "react-router-dom";
import Tooltip from "@mui/material/Tooltip";
import { formatDate, truncate } from "@/utils/textUtils";
import CircularProgress from "@mui/material/CircularProgress";
import Button from "@mui/material/Button";
import { relativeTime } from "@/utils/timeUtils";
import { useState } from "react";
import type { Blog } from "@/types/Blog";
import Grid from "@mui/material/Grid";
import Pagination from "@mui/material/Pagination";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import { Controller, useForm } from "react-hook-form";
import TextField from "@mui/material/TextField";
import RadioGroup from "@mui/material/RadioGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import Radio from "@mui/material/Radio";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { zodResolver } from "@hookform/resolvers/zod";
import { searchBlogsSchema, type SearchBlogsFormValues } from "./model/schema";
import { getTypeValidValue } from "@/utils/mapper";
import { BLOG_SORT_BY_VALUES, SORT_DIR } from "@/types/types";
import { useQueryBlogsAuthor } from "@/hooks/useQueryBlogs";
import LazyImage from "@/components/Image/LazyImage";
import placeholder from "@/assets/160x120.svg";
import IconButton from "@mui/material/IconButton";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import ViewIcon from "@mui/icons-material/AutoStories";

export default function BlogList() {
  const getInitialForm = (): SearchBlogsFormValues => ({
    title: searchParams.get("title") || "",
    content: searchParams.get("content") || "",
    sortBy: getTypeValidValue(
      searchParams.get("sortBy"),
      BLOG_SORT_BY_VALUES,
      "created_at",
    ),
    sortDir: getTypeValidValue(searchParams.get("sortDir"), SORT_DIR, "asc"),
    limit: Number.isFinite(limitParam) && limitParam > 10 ? limitParam : 15,
  });

  const [showLocaleDate, setShowLocaleDate] = useState(false);
  const [searchParams, setSearchParams] = useSearchParams();
  const limitParam = Number(searchParams.get("limit"));
  const initialForm = getInitialForm();
  const [searchForm, setSearchForm] = useState(initialForm);
  const {
    control,
    handleSubmit,
    formState: { isValid },
  } = useForm<SearchBlogsFormValues>({
    resolver: zodResolver(searchBlogsSchema),
    defaultValues: initialForm,
  });

  const [page, setPage] = useState(1);
  const onSearch = (formValues: SearchBlogsFormValues) => {
    const params = new URLSearchParams();

    Object.entries(formValues).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== "") {
        params.set(key, String(value));
      }
    });

    setSearchParams(params);
    setPage(1);
    setSearchForm(formValues);
  };

  const { data, isLoading, refetch } = useQueryBlogsAuthor(
    {
      title: searchForm.title,
      content: searchForm.content,
      sortBy: searchForm.sortBy,
      sortDir: searchForm.sortDir,
      limit: searchForm.limit,
    },
    page,
  );
  const handleRefesh = () => {
    refetch();
  };

  const handleShowCreatedAt = () => {
    setShowLocaleDate((prev) => !prev);
  };

  return (
    <Box>
      {/* Main Layout */}
      <Container sx={{ py: 4 }}>
        {/* Center */}
        <Grid size={{ xs: 12, md: 8 }} order={{ xs: 3, md: 2 }}>
          <Box component="form" onSubmit={handleSubmit(onSearch)}>
            <Stack spacing={2}>
              {/* main search */}
              <Controller
                name="title"
                control={control}
                render={({ field }) => (
                  <TextField {...field} fullWidth label="Search title" />
                )}
              />

              <Button type="submit" variant="contained" disabled={!isValid}>
                Search
              </Button>

              {/* advanced */}
              <Accordion>
                <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                  Advanced Search
                </AccordionSummary>
                <AccordionDetails>
                  <Stack spacing={2}>
                    <Controller
                      name="content"
                      control={control}
                      render={({ field }) => (
                        <TextField {...field} fullWidth label="Content" />
                      )}
                    />

                    {/* sort by */}
                    <Controller
                      name="sortBy"
                      control={control}
                      render={({ field }) => (
                        <RadioGroup row {...field}>
                          {[
                            { name: "Title", value: "title" },
                            { name: "Upload time", value: "created_at" },
                            { name: "Relevance", value: "relevance" },
                          ].map((item, index) => (
                            <FormControlLabel
                              key={index}
                              value={item.value}
                              control={<Radio />}
                              label={item.name}
                            />
                          ))}
                        </RadioGroup>
                      )}
                    />

                    {/* sort direction */}
                    <Controller
                      name="sortDir"
                      control={control}
                      render={({ field }) => (
                        <RadioGroup row {...field}>
                          {[
                            { name: "Ascending", value: "asc" },
                            { name: "Descending", value: "desc" },
                          ].map((item, index) => (
                            <FormControlLabel
                              key={index}
                              value={item.value}
                              control={<Radio />}
                              label={item.name}
                            />
                          ))}
                        </RadioGroup>
                      )}
                    />
                  </Stack>
                </AccordionDetails>
              </Accordion>
            </Stack>
          </Box>
          <Stack spacing={3} direction={{ sx: "row" }}>
            <Box
              id="action"
              sx={{
                display: "flex",
                justifyContent: "flex-end",
              }}
            >
              <Button onClick={handleRefesh}>Refesh Blogs</Button>
            </Box>
            {isLoading ? (
              <CircularProgress />
            ) : (
              (data!.blogs as Blog[])?.map((blog) => (
                <Card key={blog.blogID} sx={{ display: "flex" }}>
                  {/* Thumbnail */}
                  <Link
                    component={RouterLink}
                    to={`${blog.urlSlug}`}
                    underline="none"
                  >
                    <CardMedia>
                      <LazyImage
                        sx={{ width: "160px", minWidth: "160px" }}
                        src={blog.thumbnailUrl ?? placeholder}
                        alt="thumbnail"
                      />
                    </CardMedia>
                  </Link>

                  {/* Content */}
                  <Box
                    sx={{ display: "flex", flexDirection: "column", flex: 1 }}
                  >
                    <CardContent>
                      {/* Title */}
                      <Tooltip title={blog.title}>
                        <Typography variant="h5">
                          <Link
                            component={RouterLink}
                            to={`${blog.urlSlug}`}
                            underline="hover"
                            color="inherit"
                          >
                            {truncate(blog.title, 50)}
                          </Link>
                        </Typography>
                      </Tooltip>

                      {/* Author + Date */}
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mb: 1 }}
                      >
                        By{" "}
                        <Link
                          href={
                            "/blogs/author/" +
                            (blog.author.slug
                              ? blog.author.slug
                              : `id/${blog.author.authorID}`)
                          }
                          underline="hover"
                        >
                          {blog.author.displayName}
                        </Link>{" "}
                        •{" "}
                        <Typography
                          variant="body2"
                          component="span"
                          onClick={handleShowCreatedAt}
                        >
                          {showLocaleDate
                            ? formatDate(blog.createdAt)
                            : relativeTime(blog.createdAt)}
                        </Typography>
                      </Typography>

                      <Stack
                        direction="row"
                        spacing={1}
                        useFlexGap
                        gap={0.5}
                        flexWrap="wrap"
                        sx={{ mt: 1 }}
                      >
                        {blog.tags.map((tag) => (
                          <Typography
                            key={tag}
                            component="span"
                            sx={{
                              color: "primary.main",
                              cursor: "pointer",
                              fontWeight: 500,
                              fontSize: "0.7em",
                              "&:hover": {
                                textDecoration: "underline",
                              },
                            }}
                          >
                            <Link component={RouterLink} to={`/?tag=${tag}`}>
                              #{tag}
                            </Link>
                          </Typography>
                        ))}
                      </Stack>

                      {/* Preview */}
                      <Typography variant="body1">
                        {truncate(blog.contentText)}
                      </Typography>
                    </CardContent>
                  </Box>
                  <Stack sx={{ p: 1 }} id="action-block">
                    <Link
                      component={RouterLink}
                      to={`/author/my-blogs/${blog.urlSlug}/view`}
                    >
                      <IconButton
                        aria-label="delete"
                        color="primary"
                        title="View"
                      >
                        <ViewIcon />
                      </IconButton>
                    </Link>
                    <Link
                      component={RouterLink}
                      to={`/author/my-blogs/${blog.urlSlug}/edit`}
                    >
                      <IconButton
                        aria-label="delete"
                        color="info"
                        title="Edit blog"
                      >
                        <EditIcon />
                      </IconButton>
                    </Link>
                    <IconButton
                      aria-label="delete"
                      color="error"
                      title="Delete blog"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Stack>
                </Card>
              ))
            )}
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                mt: 1,
              }}
            >
              <Pagination
                sx={{ mt: 2 }}
                count={Math.ceil((data?.total ?? 0) / 10) || 1}
                page={page}
                onChange={(_, v) => setPage(v)}
              />
            </Box>
          </Stack>
        </Grid>
      </Container>

      {/* Footer */}
      <Box sx={{ py: 4, background: "#111" }} />
    </Box>
  );
}
