export class ActualTime{
    constructor(public id:number, public actual_time:string){}
}

export class StartEndTime{
    constructor(
        public id: number,
        public start_time: string,
        public end_time: string,
    ){}
}