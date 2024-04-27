'use client'
import Image from "next/image";
import { tableData } from "../mockdata";
import Link from "next/link";
import { tasks } from "../taskData";
import { useRouter } from "next/navigation";

export default function Task() {
  const router = useRouter();

  const tableActions = (e: React.MouseEvent<HTMLDivElement | HTMLTableRowElement>, action: string, taskId: number) => {
    e.stopPropagation();

    switch(action){
        case 'view task':
            router.push(`/dashboard/tasks/${taskId}`)
            return;
        case 'delete':
            alert('Deleting')
            return;
        case 'update':
            alert('updating task')
        default:
            return;
    }
  }
    return (
        <div className="w-full h-full flex-1 px-[3rem] py-[1.3rem] bg-white">
        <h2 className="text-[1.3rem] pb-[1rem] font-medium">List of Tasks</h2>
        <div className="w-full h-[90%] p-4 bg-gray-100 rounded-lg">
            <div className="flex justify-between items-center">
                <div className="flex justify-start items-center gap-[0.5rem]">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" className="w-5 h-5 text-red-600">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                  </svg>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5">
                    <path fill-rule="evenodd" d="M12.53 16.28a.75.75 0 0 1-1.06 0l-7.5-7.5a.75.75 0 0 1 1.06-1.06L12 14.69l6.97-6.97a.75.75 0 1 1 1.06 1.06l-7.5 7.5Z" clip-rule="evenodd" />
                  </svg>
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5">
                  <path d="M18.75 12.75h1.5a.75.75 0 0 0 0-1.5h-1.5a.75.75 0 0 0 0 1.5ZM12 6a.75.75 0 0 1 .75-.75h7.5a.75.75 0 0 1 0 1.5h-7.5A.75.75 0 0 1 12 6ZM12 18a.75.75 0 0 1 .75-.75h7.5a.75.75 0 0 1 0 1.5h-7.5A.75.75 0 0 1 12 18ZM3.75 6.75h1.5a.75.75 0 1 0 0-1.5h-1.5a.75.75 0 0 0 0 1.5ZM5.25 18.75h-1.5a.75.75 0 0 1 0-1.5h1.5a.75.75 0 0 1 0 1.5ZM3 12a.75.75 0 0 1 .75-.75h7.5a.75.75 0 0 1 0 1.5h-7.5A.75.75 0 0 1 3 12ZM9 3.75a2.25 2.25 0 1 0 0 4.5 2.25 2.25 0 0 0 0-4.5ZM12.75 12a2.25 2.25 0 1 1 4.5 0 2.25 2.25 0 0 1-4.5 0ZM9 15.75a2.25 2.25 0 1 0 0 4.5 2.25 2.25 0 0 0 0-4.5Z" />
                  </svg>
                </div>
                {/* <div className='w-fit px-[1rem] py-[0.3rem] bg-indigo-600 rounded-lg flex justify-center items-center gap-[1rem]'>
                    <div className='text-white font-semibold'>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" className="w-5 h-5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                        </svg>
                    </div>
                    <div className='text-[1rem] text-white'>
                        <Link href="/dashboard/create-task">
                            Create task
                        </Link>
                    </div>
                </div> */}
            </div>

            <div className="relative py-[2rem] overflow-x-auto">
                <table className="w-full text-left">
                  <thead >
                    <tr className="border-b border-gray-300">
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Title</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Asignee</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Status</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Hours</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {
                      tasks.map((data, index) => (
                        <tr onClick={(e) => tableActions(e, 'view task', data.id)} key={index} className="hover:bg-gray-200 cursor-pointer">
                          <td className="p-2 flex justify-start items-center gap-[1rem]">
                            <div className="font-medium ">{data.title}</div>
                          </td>
                          <td className="px-2 pb-[1.4rem] pt-[1.3rem] ">{data.asignee}</td>
                          <td className="px-2 pb-[1.4rem] pt-[1.3rem] ">{data.status}</td>
                          <td className="px-2 pb-[1.4rem] pt-[1.3rem] ">{data.hours}</td>
                          <td className="px-2 pb-[1.4rem] pt-[1.3rem] flex justify-start items-center gap-[1rem]">
                          <div onClick={e=> tableActions(e, 'delete', data.id)} className="text-red-600 font-medium">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" className="w-4 h-4">
                                  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                                </svg>
                              </div>
                              <div onClick={e=> tableActions(e, 'update', data.id)} className="p-2 bg-gray-300 text-[0.8rem] font-medium cursor-pointer rounded-md">
                                Update Task
                              </div>
                          </td>
                        </tr>
                      ))
                    }
                  </tbody>
                </table>
            </div>
        </div>
      </div>
    )
}