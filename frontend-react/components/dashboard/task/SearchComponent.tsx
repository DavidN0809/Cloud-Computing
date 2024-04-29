import React, { useEffect, useState } from 'react';
import { Box, Grid, TextField, Switch, FormControlLabel, Button, Typography, CircularProgress } from '@mui/material';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

function SearchComponent() {
  const [searchType, setSearchType] = useState('taskid');
  const [searchText, setSearchText] = useState('');
  const [searchResults, setSearchResults] = useState<TaskData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const handleSwitchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchType(event.target.checked ? 'userid' : 'taskid');
    // Reset states when search type changes
    setSearchResults([]);
    setError(null);
  };

  interface TaskData {
    id: number;
    title: string;
    description: string;
    assigned_to: string;
    status: string;
    hours: number;
    start_date: Date;
    end_date: Date;
  }
  
  const createData = (id: number, title: string, description: string, assigned_to: string, status: string, hours: number, start_date: Date, end_date: Date): TaskData => {
    return {
      id,
      title,
      description,
      assigned_to,
      status,
      hours,
      start_date,
      end_date,
    };
  };

  useEffect(() => {
    console.log("Updated searchResults:", searchResults);
  }, [searchResults]);

  const handleSearch = async () => {
    setLoading(true);
    setError(null);
    try {
      let url = searchType === 'taskid' ?
        `http://localhost:8000/tasks/get/${searchText}` :
        `http://localhost:8000/tasks/listByUser/${searchText}`;
  
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      console.log("data:", data);
  
      // If the response is an object with a 'task' property, handle it as a single task object.
      if (data && data.task) {
        const task = createData(
          data.task.id,
          data.task.title,
          data.task.description,
          data.task.assigned_to,
          data.task.status,
          data.task.hours,
          new Date(data.task.start_date),
          new Date(data.task.end_date)
        );
        setSearchResults([task]); // Set as an array with the single object
      } else {
        const tasks: TaskData[] = data.map((task: TaskData) => 
          createData(
            task.id,
            task.title,
            task.description,
            task.assigned_to,
            task.status,
            task.hours,
            new Date(task.start_date),
            new Date(task.end_date)
          )
        );
        setSearchResults(tasks);
      }
    } catch (error) {
      setError(error instanceof Error ? error : new Error('An unknown error occurred'));
    }
    setLoading(false);
  };

  return (
    <Box>
      <Box display="flex" alignItems="center" justifyContent="center" p={2}>
        <Grid container alignItems="center" spacing={2}>
          <Grid item>
            <FormControlLabel
              control={<Switch checked={searchType === 'userid'} onChange={handleSwitchChange} />}
              label={searchType === 'taskid' ? 'Search by Task ID' : 'Search by User ID'}
            />
          </Grid>
          <Grid item xs>
            <TextField
              fullWidth
              label={searchType === 'taskid' ? 'Task ID' : 'User ID'}
              variant="outlined"
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
            />
          </Grid>
          <Grid item>
            <Button onClick={handleSearch} variant="contained" color="primary" disabled={loading}>
              Search
            </Button>
          </Grid>
        </Grid>
      </Box>
      <Box p={2}>
        {loading ? (
          <CircularProgress />
        ) : error ? (
          <Typography variant="subtitle1" color="error">
            {error.message}
          </Typography>
        ) : searchResults.length > 0 ? (
          <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
              <TableHead>
                <TableRow>
                  <TableCell>Title</TableCell>
                  <TableCell align="right">ID</TableCell>
                  <TableCell align="right">Description</TableCell>
                  <TableCell align="right">Assigned To</TableCell>
                  <TableCell align="right">Status</TableCell>
                  <TableCell align="right">Hours</TableCell>
                  <TableCell align="right">Start Date</TableCell>
                  <TableCell align="right">End Date</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {searchResults.map((task) => (
                  <TableRow key={task.id}>
                    <TableCell component="th" scope="row">
                      {task.title}
                    </TableCell>
                    <TableCell align="right">{task.id}</TableCell>
                    <TableCell align="right">{task.description}</TableCell>
                    <TableCell align="right">{task.assigned_to}</TableCell>
                    <TableCell align="right">{task.status}</TableCell>
                    <TableCell align="right">{task.hours}</TableCell>
                    <TableCell align="right">{task.start_date.toDateString()}</TableCell>
                    <TableCell align="right">{task.end_date.toDateString()}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Typography variant="subtitle1" align="center">
            No results found.
          </Typography>
        )}
      </Box>
    </Box>
  );
}

export default SearchComponent;
