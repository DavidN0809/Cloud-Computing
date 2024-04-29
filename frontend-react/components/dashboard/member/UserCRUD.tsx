import React, { useState } from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import Container from '@mui/material/Container';
import Grid from '@mui/material/Grid';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import DialogActions from '@mui/material/DialogActions';
import Title from './Title';
import { useRouter } from 'next/navigation';
import { AlertColor } from '@mui/material/Alert';

type Severity = AlertColor;



export default function UserCRUD() {
  const [open, setOpen] = useState(false);
  const [openCreate, setOpenCreate] = useState(false);
  const [openUpdate, setOpenUpdate] = useState(false);
  const [openDelete, setOpenDelete] = useState(false);
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');


  const [user_id, setUser_id] = useState('');

  const router = useRouter();

  const [isAdmin, setIsAdmin] = React.useState(false);
  React.useEffect(() => {
    const cookie = document.cookie
      .split('; ')
      .find(row => row.startsWith('savedUserRole='));
    if (cookie) {
      const role = cookie.split('=')[1];
      setIsAdmin(role === 'admin');
    }
  }, []);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setOpenCreate(false);
    setOpenUpdate(false);
    setOpenDelete(false);
  };

  const handleClickOpenCreate = () => {
    setOpenCreate(true);
  };

  const handleClickOpenUpdate = () => {
    setOpenUpdate(true);
  };

  const handleClickOpenDelete = () => {
    setOpenDelete(true);
  };
  

  const handleCreate = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const cookies = document.cookie;
    const cookieName = 'token=';
    const cookieArray = cookies.split('; '); // cookies are separated by '; '
    // Find the cookie value for 'savedUserRole'
    const tokenCookie = cookieArray.find(cookie => cookie.startsWith(cookieName));
    const admin_token = tokenCookie ? tokenCookie.split('=')[1] : null;
    type TaskData = {
      username: string;
      email: string;
      password: string;
    };
    
    const taskData: TaskData = {
      username: username,
      email: email,
      password: password,
      // 现在你可以条件性地添加 parent_task
    };
    

    
    // const testData = {
    //   "title":"5",
    //   "description":"5",
    //   "assigned_to":"assignedTo",
    //   "status":"plan",
    //   "hours":8,
    //   "start_date":"2024-04-28T00:00:00.000Z",
    //   "end_date":"2024-04-30T00:00:00.000Z"
    // };

    try {
      const response = await fetch('http://localhost:8000/users/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(taskData)
      });
      console.log(JSON.stringify(taskData));
      console.log("response.ok:",response.ok);
      console.log("response.text:",response.text);
      console.log("response:",response);
      const data = await response.json(); // Assuming the server responds with JSON
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      console.log('user created', data);
      window.location.href = '/dashboard/member?stat=succeed';
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      window.location.href = '/dashboard/member?stat=failed';
    } finally {
    handleClose(); // This will close the dialog in any case after operation
  }
  };


  const handleUpdate = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const cookies = document.cookie;
    const cookieName = 'token=';
    const cookieArray = cookies.split('; '); // cookies are separated by '; '
    // Find the cookie value for 'savedUserRole'
    const tokenCookie = cookieArray.find(cookie => cookie.startsWith(cookieName));
    const admin_token = tokenCookie ? tokenCookie.split('=')[1] : null;

    type TaskData = {
      username: string;
      email: string;
      password: string;
    };
    
    const taskData: TaskData = {
      username: username,
      email: email,
      password: password,
      // 现在你可以条件性地添加 parent_task
    };

    try {
      const response = await fetch(`http://localhost:8000/users/update/${user_id}`, {
        method: 'PUT', 
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${admin_token}`  
        },
        body: JSON.stringify(taskData)
      });
      console.log(`http://localhost:8000/users/update/${user_id}`);
      console.log(JSON.stringify(taskData));
      console.log("response.ok:",response.ok);
      console.log("response.text:",response.text);
      console.log("response:",response);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      console.log('user updated');
      window.location.href = '/dashboard/member?stat=succeed';
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      window.location.href = '/dashboard/member?stat=failed';
    } finally {
    handleClose(); // This will close the dialog in any case after operation
  }
  };


  const handleDelete = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();

    

    try {
      const cookies = document.cookie;
      const cookieName = 'token=';
      const cookieArray = cookies.split('; '); // cookies are separated by '; '

      // Find the cookie value for 'savedUserRole'
      const tokenCookie = cookieArray.find(cookie => cookie.startsWith(cookieName));
      const admin_token = tokenCookie ? tokenCookie.split('=')[1] : null;

      const response = await fetch(`http://localhost:8000/users/remove/${user_id}`, {
        method: 'DELETE', 
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${admin_token}`  
        },
      });
      console.log(`http://localhost:8000/users/remove/${user_id}`);
      console.log("response.ok:",response.ok);
      console.log("response.text:",response.text);
      console.log("response:",response);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      console.log('user removed');
      window.location.href = '/dashboard/member?stat=succeed';
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      window.location.href = '/dashboard/member?stat=failed';
      
    } finally {
    handleClose(); // This will close the dialog in any case after operation
  }
  };

  return (
    <React.Fragment>
      <Title>Today</Title>
      <Container maxWidth="lg">
        <Grid container spacing={3} alignItems="center" justifyContent="space-around">
          <Grid item md={4} display="flex" justifyContent="center">
            <Button variant="outlined" color="primary" onClick={handleClickOpenCreate}>
              Create User
            </Button>
            <Dialog open={openCreate} onClose={handleClose}>
              <DialogTitle>Create New User</DialogTitle>
              <DialogContent>
                <TextField autoFocus margin="dense" id="username" label="User name" type="text" fullWidth variant="standard" value={username} onChange={(e) => setUsername(e.target.value)} />
                <TextField margin="dense" id="email" label="email" type="text" fullWidth multiline variant="standard" value={email} onChange={(e) => setEmail(e.target.value)} />
                <TextField margin="dense" id="password" label="password" type="number" fullWidth variant="standard" value={password} onChange={(e) => setPassword(e.target.value)} />
              </DialogContent>
              <DialogActions>
                <Button onClick={handleClose} color="secondary">
                  Cancel
                </Button>
                <Button onClick={handleCreate} color="primary">
                  Create
                </Button>
              </DialogActions>
            </Dialog>
          </Grid>

          <Grid item md={4} display="flex" justifyContent="center">
            <Button variant="outlined" color="primary" onClick={handleClickOpenUpdate}>
              Update User
            </Button>
            <Dialog open={openUpdate} onClose={handleClose}>
              <DialogTitle>Upadte User</DialogTitle>
              <DialogContent>
                
                <TextField autoFocus margin="dense" id="username" label="User name" type="text" fullWidth variant="standard" value={username} onChange={(e) => setUsername(e.target.value)} />
                <TextField margin="dense" id="email" label="email" type="text" fullWidth multiline variant="standard" value={email} onChange={(e) => setEmail(e.target.value)} />
                <TextField margin="dense" id="password" label="password" type="number" fullWidth variant="standard" value={password} onChange={(e) => setPassword(e.target.value)} />
                <TextField margin="dense" id="user_id" label="User id" type="text" fullWidth variant="standard" value={user_id} onChange={(e) => setUser_id(e.target.value)} />
              </DialogContent>
              <DialogActions>
                <Button onClick={handleClose} color="secondary">
                  Cancel
                </Button>
                <Button onClick={handleUpdate} color="primary">
                  Update
                </Button>
              </DialogActions>
            </Dialog>
          </Grid>


          <Grid item md={4} display="flex" justifyContent="center">
          {isAdmin ? (
            <>
              <Button variant="outlined" color="primary" onClick={handleClickOpenDelete}>
                Delete User
              </Button>
              <Dialog open={openDelete} onClose={handleClose}>
                <DialogTitle>Delete User</DialogTitle>
                <DialogContent>
                <TextField autoFocus margin="dense" id="user_id" label="User id" type="text" fullWidth variant="standard" value={user_id} onChange={(e) => setUser_id(e.target.value)} />
                </DialogContent>
                <DialogActions>
                  <Button onClick={handleClose} color="secondary">
                    Cancel
                  </Button>
                  <Button onClick={handleDelete} color="primary">
                    Delete
                  </Button>
                </DialogActions>
              </Dialog>
            </>
          ) : (
            <Button variant="outlined" color="secondary">
              Access Denied
            </Button>
          )}
        </Grid>







        </Grid>
      </Container>
    </React.Fragment>
  );
}
