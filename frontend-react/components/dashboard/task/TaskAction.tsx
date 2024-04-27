import * as React from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import DialogActions from '@mui/material/DialogActions';
import Title from './Title';

export default function TaskAction({ returnState }: { returnState: boolean }) {
  const [open, setOpen] = React.useState(false);
  const [title, setTitle] = React.useState('');
  const [description, setDescription] = React.useState('');
  const [assignedTo, setAssignedTo] = React.useState('');
  const [status, setStatus] = React.useState('');
  const [hours, setHours] = React.useState('');

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleCreate = async () => {
    const taskData = {
      title: title,
      description: description,
      assigned_to: assignedTo,
      status: status,
      hours: hours
    };
    
    try {
      const response = await fetch('http://localhost:8000/tasks/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(taskData)
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      } else {
        console.log('Task created', await response.json());
        setOpen(false);
        returnState(true); // Call the callback with success status
      }
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      returnState(false); // Call the callback with failure status
    }
  };

  return (
    <React.Fragment>
      <Title>Today</Title>
      <div>
        <Button variant="outlined" color="primary" onClick={handleClickOpen}>
          Create Task
        </Button>
        <Dialog open={open} onClose={handleClose}>
          <DialogTitle>Create New Task</DialogTitle>
          <DialogContent>
            <TextField autoFocus margin="dense" id="title" label="Task Name" type="text" fullWidth variant="standard" value={title} onChange={(e) => setTitle(e.target.value)} />
            <TextField margin="dense" id="description" label="Task Description" type="text" fullWidth multiline variant="standard" value={description} onChange={(e) => setDescription(e.target.value)} />
            <TextField margin="dense" id="assigned_to" label="Assigned To" type="text" fullWidth variant="standard" value={assignedTo} onChange={(e) => setAssignedTo(e.target.value)} />
            <TextField margin="dense" id="status" label="Status" type="text" fullWidth variant="standard" value={status} onChange={(e) => setStatus(e.target.value)} />
            <TextField margin="dense" id="hours" label="Hours" type="number" fullWidth variant="standard" value={hours} onChange={(e) => setHours(e.target.value)} />
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
      </div>
    </React.Fragment>
  );
}
