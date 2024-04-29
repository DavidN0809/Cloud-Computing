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



export default function TaskAction() {
  const [open, setOpen] = useState(false);
  const [openCreate, setOpenCreate] = useState(false);
  const [openUpdate, setOpenUpdate] = useState(false);
  const [openDelete, setOpenDelete] = useState(false);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [assignedTo, setAssignedTo] = useState('');
  const [status, setStatus] = useState('');
  const [hours, setHours] = useState('');
  const [start_date, setStartDate] = useState('');
  const [end_date, setEndDate] = useState('');
  const [task_id, setTaskId] = useState('');
  const [parent_task, setParentTask] = useState('');
  const router = useRouter();
  const [message, setMessage] = React.useState('');
  const [severity, setSeverity] = React.useState<Severity>('success');

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
    const isoStartDate = new Date(start_date).toISOString();   //prepare data send to api
    const isoEndDate = new Date(end_date).toISOString();
    type TaskData = {
      title: string;
      description: string;
      assigned_to: string;
      status: string;
      hours: number;
      start_date: string;
      end_date: string;
      parent_task?: string; // 这行添加了一个可选属性
    };
    
    const taskData: TaskData = {
      title: title,
      description: description,
      assigned_to: assignedTo,
      status: status,
      hours: parseInt(hours, 10),
      start_date: isoStartDate,
      end_date: isoEndDate,
      // 现在你可以条件性地添加 parent_task
    };
    
    // Conditionally add the parent_task if it is not empty
    if (parent_task !== "") {
      taskData.parent_task = parent_task;
    }

    
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
      const response = await fetch('http://localhost:8000/tasks/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
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
  
      console.log('Task created', data);
      //window.location.href = '/dashboard/tasks?stat=succeed';
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      //window.location.href = '/dashboard/tasks?stat=failed';
    } finally {
    handleClose(); // This will close the dialog in any case after operation
  }
  };


  const handleUpdate = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const isoStartDate = new Date(start_date).toISOString();   //prepare data send to api
    const isoEndDate = new Date(end_date).toISOString();
    

    const taskData = {
      "title":title,
      "description":description,
      "assigned_to":assignedTo,
      "status":status,
      "hours":parseInt(hours, 10),
      "start_date":isoStartDate,
      "end_date":isoEndDate
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
      window.location.href = '/dashboard/tasks?stat=succeed';
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      window.location.href = '/dashboard/tasks?stat=failed';
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

      const response = await fetch(`http://localhost:8000/tasks/remove/${task_id}`, {
        method: 'DELETE', 
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${admin_token}`  
        },
      });
      console.log(`http://localhost:8000/tasks/remove/${task_id}`);
      console.log("response.ok:",response.ok);
      console.log("response.text:",response.text);
      console.log("response:",response);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      console.log('Task removed');
      window.location.href = '/dashboard/tasks?stat=succeed';
      
      
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
      window.location.href = '/dashboard/tasks?stat=failed';
      
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
              Create Task
            </Button>
            <Dialog open={openCreate} onClose={handleClose}>
              <DialogTitle>Create New Task</DialogTitle>
              <DialogContent>
                <TextField autoFocus margin="dense" id="title" label="Task Name" type="text" fullWidth variant="standard" value={title} onChange={(e) => setTitle(e.target.value)} />
                <TextField margin="dense" id="description" label="Task Description" type="text" fullWidth multiline variant="standard" value={description} onChange={(e) => setDescription(e.target.value)} />
                <TextField margin="dense" id="assigned_to" label="Assigned To" type="text" fullWidth variant="standard" value={assignedTo} onChange={(e) => setAssignedTo(e.target.value)} />
                <TextField margin="dense" id="status" label="Status" type="text" fullWidth variant="standard" value={status} onChange={(e) => setStatus(e.target.value)} />
                <TextField margin="dense" id="hours" label="Hours" type="number" fullWidth variant="standard" value={hours} onChange={(e) => setHours(e.target.value)} />
                <TextField margin="dense" id="start_date" label="start date" type="date" fullWidth variant="standard" value={start_date} onChange={(e) => setStartDate(e.target.value)} />
                <TextField margin="dense" id="end_date" label="end date" type="date" fullWidth variant="standard" value={end_date} onChange={(e) => setEndDate(e.target.value)} />
                <TextField margin="dense" id="parent_task" label="parent task date" type="text" fullWidth variant="standard" value={parent_task} onChange={(e) => setParentTask(e.target.value)} />
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
              Update Task
            </Button>
            <Dialog open={openUpdate} onClose={handleClose}>
              <DialogTitle>Upadte Task</DialogTitle>
              <DialogContent>
                <TextField autoFocus margin="dense" id="task_id" label="Task ID" type="text" fullWidth variant="standard" value={task_id} onChange={(e) => setTaskId(e.target.value)} />
                <TextField margin="dense" id="title" label="Task Name" type="text" fullWidth variant="standard" value={title} onChange={(e) => setTitle(e.target.value)} />
                <TextField margin="dense" id="description" label="Task Description" type="text" fullWidth multiline variant="standard" value={description} onChange={(e) => setDescription(e.target.value)} />
                <TextField margin="dense" id="assigned_to" label="Assigned To" type="text" fullWidth variant="standard" value={assignedTo} onChange={(e) => setAssignedTo(e.target.value)} />
                <TextField margin="dense" id="status" label="Status" type="text" fullWidth variant="standard" value={status} onChange={(e) => setStatus(e.target.value)} />
                <TextField margin="dense" id="hours" label="Hours" type="number" fullWidth variant="standard" value={hours} onChange={(e) => setHours(e.target.value)} />
                <TextField margin="dense" id="start_date" label="start date" type="date" fullWidth variant="standard" value={start_date} onChange={(e) => setStartDate(e.target.value)} />
                <TextField margin="dense" id="end_date" label="end date" type="date" fullWidth variant="standard" value={end_date} onChange={(e) => setEndDate(e.target.value)} />
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
                Delete Task
              </Button>
              <Dialog open={openDelete} onClose={handleClose}>
                <DialogTitle>Delete Task</DialogTitle>
                <DialogContent>
                  <TextField autoFocus margin="dense" id="task_id" label="Task ID" type="text" fullWidth variant="standard" value={task_id} onChange={(e) => setTaskId(e.target.value)} />
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
