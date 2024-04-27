'use client';
import Image from "next/image";
import { tableData } from "./mockdata";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function Page() {
  const router = useRouter();

  const tableActions = (e: React.MouseEvent<HTMLDivElement | HTMLTableRowElement>, type: string, userId: number) => {
    e.stopPropagation();

    switch(type){
        case 'view table':
            router.push(`/dashboard/${userId}`)
            return;
        case 'delete':
            alert('Deleting user')
            return;
        case 'update':
            alert('updating user')
            return;
        case 'assign task':
          router.push(`/dashboard/create-task/${userId}`)
            return;
        default:
            return;
    }
}
    return (
      <div className="w-full h-full flex-1 px-[2.5rem] py-[1.3rem] bg-white">
        <h2 className="text-[1.3rem] pb-[1rem] font-medium">Role administration</h2>
        <div className="w-full h-[90%] py-[1.5rem] bg-gray-100 rounded-lg">
          <div className="px-[1.6rem] flex justify-start items-center gap-[0.5rem]">
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

            <div className="relative py-[2rem] overflow-x-auto">
                <table className="w-full text-left">
                  <thead >
                    <tr className="px-[2rem] border-b border-gray-300">
                      <th className="pl-[2rem] pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Username</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Email</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Login</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Role</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Actions</th>
                      <th className="px-2 pt-2 pb-[0.8rem] font-medium text-[1.1rem]">Work</th>
                    </tr>
                  </thead>
                  <tbody>
                    {
                      tableData.map((data, index) => (
                        // <Link className="w-full" href={`/dashboard/${data.id}`}>
                          <tr onClick={(e) => tableActions(e, 'view table', data.id)} key={index} className="w-full hover:bg-gray-200 hover:scale-105 transition-all duration-100 ease-linear cursor-pointer">
                            <td className="pl-[1.5rem] p-2 flex justify-start items-center gap-[1rem]">
                              
                                {/* <Image src="/img/profile.jpeg" alt="Profile" width={50} height={50} className="rounded-full"/> */}
                              <div className="font-medium ">{data.name}</div>
                            </td>
                            <td className="px-2 pb-[1.4rem] pt-[1.3rem] ">{data.email}</td>
                            <td className="px-2 pb-[1.4rem] pt-[1.3rem] ">{data.login}</td>
                            <td className="px-2 pb-[1.4rem] pt-[1.3rem] ">{data.role}</td>
                            <td className="px-2 pb-[1.4rem] pt-[1.3rem] flex justify-start items-center gap-[1rem]">
                              <div onClick={e=> tableActions(e, 'delete', data.id)} className="text-red-600 font-medium">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" className="w-4 h-4">
                                  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                                </svg>
                              </div>
                              <div onClick={e=> tableActions(e, 'update', data.id)} className="p-2 bg-gray-300 text-[0.8rem] font-medium cursor-pointer rounded-md">
                                Update user
                              </div>
                            </td>
                            <td className="pr-[1.5rem]">
                              <div onClick={(e) => tableActions(e, 'assign task', data.id)} className="p-[5px] text-white bg-indigo-500 rounded-md text-center">Assign task</div>
                            </td>
                          </tr>
                        // </Link>
                      ))
                    }
                  </tbody>
                </table>
            </div>
        </div>
      </div>
    )
  }