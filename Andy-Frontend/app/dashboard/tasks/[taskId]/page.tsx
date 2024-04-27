import Link from "next/link"

export default function TaskDetails({params}: {
    params: {taskId: string}
}) {
    return (
        <div className="w-full py-[3rem] bg-white">
            <div className="w-[60%] mx-auto py-4 px-[3rem] bg-gray-100 rounded-lg">
                <h2 className="py-[1rem] text-[1.3rem] font-medium border-b border-gray-300">Task Details</h2>
                <div className="w-full py-4 flex justify-end items-center gap-[1rem] ">
                    <div className="px-[0.6rem] py-[0.4rem] flex justify-center items-center gap-[1rem] text-red-600 font-medium bg-gray-300 rounded-lg cursor-pointer">
                        <div>Delete Task</div>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" className="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                        </svg>
                    </div>
                    <Link href="/dashboard/create-user" className="px-[0.6rem] py-[0.4rem] flex justify-center items-center gap-[1rem] text-indigo-600 font-medium bg-gray-300 rounded-lg">
                        <div>Update Task</div>
                    </Link>
                    
                </div>
                <div className="pb-3 flex justify-start items-center">
                    <div className="pr-[3rem] text-[1rem]">Title: </div>
                    <div className="text-[1.1rem] font-medium">Michael Rhule</div>
                </div>
                <div className="pb-3 flex justify-start items-center">
                    <div className="pr-[3rem] text-[1rem]">Description: </div>
                    <div className="text-[1.1rem] font-medium">rhule@gmail.com</div>
                </div>
                <div className="pb-3 flex justify-start items-center">
                    <div className="pr-[3rem] text-[1rem]">Asignee: </div>
                    <div className="text-[1.1rem] font-medium">lksjdfie83894</div>
                </div>
                <div className="pb-3 flex justify-start items-center">
                    <div className="pr-[3rem] text-[1rem]">Status: </div>
                    <div className="text-[1.1rem] font-medium">Pending</div>
                </div>
                <div className="pb-3 flex justify-start items-center">
                    <div className="pr-[3rem] text-[1rem]">Hours: </div>
                    <div className="text-[1.1rem] font-medium">23</div>
                </div>
            </div>
        </div>
    )
}