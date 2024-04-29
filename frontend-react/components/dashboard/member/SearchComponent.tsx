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
  const [searchType, setSearchType] = useState('userid');
  const [searchText, setSearchText] = useState('');
  const [searchResults, setSearchResults] = useState<UserData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const handleSwitchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchType(event.target.checked ? 'userid' : 'userid'); // Adjust if needed for different user search types
    // Reset states when search type changes
    setSearchResults([]);
    setError(null);
  };

  interface UserData {
    id: number;
    username: string;
    email: string;
  }
  
  const createData = (id: number, username: string, email: string): UserData => {
    return { id, username, email };
  };

  useEffect(() => {
    console.log("Updated searchResults:", searchResults);
  }, [searchResults]);

  const handleSearch = async () => {
    setLoading(true);
    setError(null);
    try {
      const url = `http://localhost:8000/users/get/${searchText}`;

      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      console.log("data:", data);
  
      // Assume the response is an array of users or a single user object
      if (Array.isArray(data)) {
        const users = data.map((user) => createData(user.id, user.username, user.email));
        setSearchResults(users);
      } else {
        const user = createData(data.id, data.username, data.email);
        setSearchResults([user]); // Set as an array with the single object
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
              label={searchType === 'userid' ? 'Search by User ID' : 'Search by User ID'}
            />
          </Grid>
          <Grid item xs>
            <TextField
              fullWidth
              label="User ID"
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
                  <TableCell>ID</TableCell>
                  <TableCell align="right">Username</TableCell>
                  <TableCell align="right">Email</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {searchResults.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell component="th" scope="row">
                      {user.id}
                    </TableCell>
                    <TableCell align="right">{user.username}</TableCell>
                    <TableCell align="right">{user.email}</TableCell>
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
``
