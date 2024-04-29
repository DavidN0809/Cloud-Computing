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
  const [searchResults, setSearchResults] = useState<BillingData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const handleSwitchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchType(event.target.checked ? 'userid' : 'taskid');
    // Reset states when search type changes
    setSearchResults([]);
    setError(null);
  };

  interface BillingData {
    user_id: string;
    task_id: number;
    hours: number;
    hourly_rate: number;
    amount: number;
  }
  
  function createData(
    user_id: string,
    task_id: number,
    hours: number,
    hourly_rate: number,
    amount: number,
  ): BillingData {
    return {
      user_id,
      task_id,
      hours,
      hourly_rate,
      amount,
    };
  }

  const cookies = document.cookie;
  const cookieName = 'token=';
  const cookieArray = cookies.split('; '); 
  const tokenCookie = cookieArray.find(cookie => cookie.startsWith(cookieName));
  const admin_token = tokenCookie ? tokenCookie.split('=')[1] : null;

  useEffect(() => {
    console.log("Updated searchResults:", searchResults);
  }, [searchResults]);

  const handleSearch = async () => {
    setLoading(true);
    setError(null);
    try {
      let url = searchType === 'taskid' ?
        `http://localhost:8000/billings/get/${searchText}` :
        `http://localhost:8000/billings/list?user_id=${searchText}`;
  
      const response = await fetch(url,{
        headers: {
          'Authorization': `Bearer ${admin_token}`  
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      console.log("data:", data);
  
      // If the response is an object with a 'task' property, handle it as a single task object.
      if (searchType === 'taskid') {
        const task = createData(
          data.user_id,
          data.task_id,
          data.hours,
          data.hourly_rate,
          data.amount,
        );
        setSearchResults([task]); // Set as an array with the single object
      } else {
        const tasks: BillingData[] = data.map((task: BillingData) => 
          createData(
            task.user_id,
            task.task_id,
            task.hours,
            task.hourly_rate,
            task.amount,
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
              label={searchType === 'taskid' ? 'Search by Billing ID' : 'Search by User ID'}
            />
          </Grid>
          <Grid item xs>
            <TextField
              fullWidth
              label={searchType === 'taskid' ? 'Billing ID' : 'User ID'}
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
                  <TableCell>user_id</TableCell>
                  <TableCell align="right">task_id</TableCell>
                  <TableCell align="right">hours</TableCell>
                  <TableCell align="right">hourly rate</TableCell>
                  <TableCell align="right">amount</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {searchResults.map((task) => (
                  <TableRow key={task.task_id}>
                    <TableCell component="th" scope="row">
                      {task.user_id}
                    </TableCell>
                    <TableCell align="right">{task.task_id}</TableCell>
                    <TableCell align="right">{task.hours}</TableCell>
                    <TableCell align="right">{task.hourly_rate}</TableCell>
                    <TableCell align="right">{task.amount}</TableCell>
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
