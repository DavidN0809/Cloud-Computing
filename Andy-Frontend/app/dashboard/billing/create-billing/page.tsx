import Modal from "@/components/Modal";

export default function CreateBill({params}: {
    params: {userId: string}
}) {
    return (
        <Modal>
            <div className="w-[550px] h-[500px] py-[3rem] bg-white rounded-lg">
                <h2 className="w-fit pb-[3rem] mx-auto text-[1.3rem] font-semibold uppercase">Billing Form</h2>
                <form className="max-w-md mx-auto ">
                    <div className="mb-[2rem]">
                        <label  className="block mb-2 text-md font-medium text-gray-900 dark:text-white">Hours</label>
                        <input type="number" id="hours" className="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 dark:shadow-sm-light" placeholder="Coding staff" required />
                    </div>
                    <div className="mb-[1.8rem]">
                        <label className="block mb-2 text-md font-medium text-gray-900 dark:text-white">Amount</label>
                        <input type="number" id="amount" className="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 dark:shadow-sm-light" placeholder="codingstaff@gmail.com" required />
                    </div>
                    <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Create a bill</button>
                </form>

            </div>
        </Modal>
    )
}