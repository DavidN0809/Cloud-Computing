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
import { useRouter } from 'next/router';
import { AlertColor } from '@mui/material/Alert';

type Severity = AlertColor;

export default function UserCRUD() {
  const [open, setOpen] = useState(false);
  const [userId, setUserId] = useState('');
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();
  const [severity, setSeverity] = useState<Severity>('success');
  const [message, setMessage] = useState('');
  const [mode, setMode] = useState('create'); // "create", "update", "delete"

  const handleOpen = (newMode: string, id = '') => {
    setMode(newMode);
    setUserId(id);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleCreate = async () => {
    try {
      const userData = { username, email, password };
      const response = await fetch('http://localhost:8000/users/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData)
      });

      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      await response.json();
      setMessage('User created successfully');
      setSeverity('success');
      router.push('/users/list');
    } catch (error) {
      console.error('Error creating user:', error);
      setMessage('Error creating user');
      setSeverity('error');
      router.push('/users/list?status=failed');
    } finally {
      handleClose();
    }
  };

  const handleUpdate = async () => {
    try {
      const userData = { username, email, password };
      const response = await fetch(`http://localhost:8000/users/update/${userId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData)
      });

      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      await response.json();
      setMessage('User updated successfully');
      setSeverity('success');
      router.push('/users/list');
    } catch (error) {
      console.error('Error updating user:', error);
      setMessage('Error updating user');
      setSeverity('error');
      router.push('/users/list?status=failed');
    } finally {
      handleClose();
    }
  };

  const handleDelete = async () => {
    try {
      const response = await fetch(`http://localhost:8000/users/remove/${userId}`, {
        method: 'DELETE'
      });

      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      await response.json();
      setMessage('User deleted successfully');
      setSeverity('success');
      router.push('/users/list');
    } catch (error) {
      console.error('Error deleting user:', error);
      setMessage('Error deleting user');
      setSeverity('error');
      router.push('/users/list?status=failed');
    } finally {
      handleClose();
    }
  };

  return (
    <React.Fragment>
      <Title>User Management</Title>
      <Container maxWidth="lg">
        <Grid container spacing={2} alignItems="center" justifyContent="space-around">
          <Grid item md={6} display="flex" justifyContent="center">
            <Button variant="outlined" color="primary" onClick={() => handleOpen('create')}>
              Create User
            </Button>
            <Button variant="outlined" color="primary" onClick={() => handleOpen('update', userId)}>
              Update User
            </Button>
            <Button variant="outlined" color="primary" onClick={() => handleOpen('delete', userId)}>
              Delete User
            </Button>
          </Grid>
        </Grid>
        <Dialog open={open} onClose={handleClose}>
          <DialogTitle>{`${mode.charAt(0).toUpperCase() + mode.slice(1)} User`}</DialogTitle>
          <DialogContent>
            {mode !== 'delete' && (
              <>
                <TextField autoFocus margin="dense" id="username" label="Username" type="text" fullWidth variant="standard" value={username} onChange={(e) => setUsername(e.target.value)} />
                <TextField margin="dense" id="email" label="Email" type="email" fullWidth variant="standard" value={email} onChange={(e) => setEmail(e.target.value)} />
                <TextField margin="dense" id="password" label="Password" type="password" fullWidth variant="standard" value={password} onChange={(e) => setPassword(e.target.value)} />
              </>
            )}
            {mode === 'delete' && <p>Are you sure you want to delete this user?</p>}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="secondary">
              Cancel
            </Button>
            {mode === 'create' && <Button onClick={handleCreate} color="primary">Create</Button>}
            {mode === 'update' && <Button onClick={handleUpdate} color="primary">Update</Button>}
            {mode === 'delete' && <Button onClick={handleDelete} color="primary">Delete</Button>}
          </DialogActions>
        </Dialog>
      </Container>
    </React.Fragment>
  );
}
