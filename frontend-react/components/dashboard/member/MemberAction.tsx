import React, { useState } from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import Container from '@mui/material/Container';
import Grid from '@mui/material/Grid';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import DialogActions from '@mui/material/DialogActions';
import Title from '../Title';
import { useRouter } from 'next/navigation';
import { AlertColor } from '@mui/material/Alert';

type Severity = AlertColor;



export default function TaskAction() {
  const [open, setOpen] = useState(false);
  const [openCreate, setOpenCreate] = useState(false);
  const [openUpdate, setOpenUpdate] = useState(false);
  const [username, setTitle] = useState('');
  const [email, setDescription] = useState('');
  const [password, setpassword] = useState('');
  const [task_id, setTaskId] = useState('');
  const router = useRouter();
  const [message, setMessage] = React.useState('');
  const [severity, setSeverity] = React.useState<Severity>('success');

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setOpenCreate(false);
    setOpenUpdate(false);
  };

  const handleClickOpenCreate = () => {
    setOpenCreate(true);
  };

  const handleClickOpenUpdate = () => {
    setOpenUpdate(true);
  };
  

  const handleCreate = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const userData = {
      "username":username,
      "email":email,
      "password":password,
    };
    try {
      const response = await fetch('http://localhost:8000/users/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
      });
      console.log(JSON.stringify(userData));
      console.log("response.ok:",response.ok);
      console.log("response.text:",response.text);
      console.log("response:",response);
      const data = await response.json(); // Assuming the server responds with JSON
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      console.log('User created', data);
      await router.push('/dashboard/members?stat=succeed');
      refreshPage();

      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      router.push('/dashboard/members?stat=failed');
    } finally {
    handleClose(); // This will close the dialog in any case after operation
  }
  };

  const refreshPage = () => {
    // Reload the page
    window.location.reload();
  };

  const handleUpdate = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const taskData = {
      "username":username,
      "email":email,
      "password":password,
      
    };
    try {
      const response = await fetch(`http://localhost:8000/tasks/update/${task_id}`, {
        method: 'PUT', 
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(taskData)
      });
      console.log(`http://localhost:8000/tasks/update/${task_id}`);
      console.log(JSON.stringify(taskData));
      console.log("response.ok:",response.ok);
      console.log("response.text:",response.text);
      console.log("response:",response);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      console.log('Task updated');
      router.push('/dashboard/tasks?stat=succeed');
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      router.push('/dashboard/tasks?stat=failed');
    } finally {
    handleClose(); // This will close the dialog in any case after operation
  }
  };

  return (
    <React.Fragment>
      <Title>Add User</Title>
      <Container maxWidth="lg">
        <Grid container spacing={2} alignItems="center" justifyContent="space-around">
          <Grid item md={6} display="flex" justifyContent="center">
            <Button variant="outlined" color="primary" onClick={handleClickOpenCreate}>
              Add User
            </Button>
            <Dialog open={openCreate} onClose={handleClose}>
              <DialogTitle>Add New User</DialogTitle>
              <DialogContent>
                <TextField autoFocus margin="dense" id="username" label="Username" type="text" required fullWidth variant="standard" value={username} onChange={(e) => setTitle(e.target.value)} />
                <TextField margin="dense" id="email" label="Email" type="text" fullWidth multiline variant="standard" value={email} onChange={(e) => setDescription(e.target.value)} />
                <TextField margin="dense" id="password" label="Password" type="password" fullWidth variant="standard" value={password} onChange={(e) => setpassword(e.target.value)} />
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

          <Grid item md={6} display="flex" justifyContent="center">
            <Button variant="outlined" color="primary" onClick={handleClickOpenUpdate}>
              Update User
            </Button>
            <Dialog open={openUpdate} onClose={handleClose}>
              <DialogTitle>Upadte User</DialogTitle>
              <DialogContent>
                <TextField autoFocus margin="dense" id="username" label="Username" type="text" fullWidth variant="standard" value={task_id} onChange={(e) => setTaskId(e.target.value)} />
                <TextField margin="dense" id="email" label="Email" type="text" fullWidth variant="standard" value={username} onChange={(e) => setTitle(e.target.value)} />
                <TextField margin="dense" id="password" label="Password" type="text" fullWidth multiline variant="standard" value={email} onChange={(e) => setDescription(e.target.value)} />
                
               
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
        </Grid>
      </Container>
    </React.Fragment>
  );
}
