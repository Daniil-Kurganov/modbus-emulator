import { Component } from "@angular/core";
import {FormsModule} from "@angular/forms";
import { HttpClient, HttpClientModule, HttpHandler } from "@angular/common/http";
import { interval, Subscription } from 'rxjs';
import { ActualTime, StartEndTime } from "./times"
     
@Component({
    selector: "modbus-emulator",
    standalone: true,
    imports: [FormsModule, HttpClientModule],
    templateUrl: './app.component.html'
})

export class AppComponent{  
    actual_times: ActualTime[]=[];
    server_id: number = -1
    status: string = "Waiting"
    rewind_time: string
    constructor(private http: HttpClient) {}
    URLHead = "modbus-emulator" 
    getActualTime(): void {
        let request = `${this.URLHead}/time/actual`
        if (this.server_id!==-1) {
            request = request + `?server_id=${this.server_id}`
        }
        if (this.actual_times.length !== 0 && this.server_id > this.actual_times.length) {
            console.error("Error: server ID must be in [0:%d]", (this.actual_times.length + 1))
            this.status = `Error - server ID must be in [0:${this.actual_times.length + 1}]`
            return
        }
        this.http.get(request).subscribe({next:(data:any) => this.actual_times=data, error: error => console.log(error)});
        this.status = "Success"
    }
    rewindEmulation(): void {
        this.rewind_time = this.rewind_time.replace("T", "%20")
        this.rewind_time = this.rewind_time.replaceAll(":", "%3A")
        console.log(this.rewind_time)
        let request = `${this.URLHead}/time/rewind_emulation?timepoint=${this.rewind_time}`
        if (this.server_id!==-1) {
            request = request + `&server_id=${this.server_id}`
        }
        this.http.post(request, {}).subscribe({error: error => console.log(error)});
        this.status = "Success"
    }
    slaveID: number
    registerType: string
    operation: string
    startIndex: string
    count: number
    registerData: string
    manuallyRegistersWork(): void {
        let isDecimalCase = true
        if (this.startIndex.substring(0, 2) === "0x") {isDecimalCase = false}
        let request = `${this.URLHead}/registers?server_id=${this.server_id}&slave_id=${this.slaveID}&type=${this.registerType}&start_index=${parseInt(this.startIndex)}`
        if (this.operation === "read") {
            request += `&count=${this.count}`
            this.http.get(request).subscribe({next:(data:any) => this.registerData=data, error: error => console.error(error)});
            return
        } else {
            let newRegistersData: number[]=[];
            let registerDataArray = this.registerData.split(" ")
            for (let index = 0; index < registerDataArray.length; index++) {
                const element = registerDataArray[index];
                if (!isDecimalCase) {
                    if (element.substring(0, 2) !== "0x") {
                        this.status = `Error: current number mode is hexadecimal, but current register's data is ${element}`
                        return
                    }
                }
                newRegistersData.push(parseInt(element))
            }
            this.http.post(request, {"registers": newRegistersData}).subscribe({error: error => console.log(error)});
        }
    }
    subscription = interval(500).subscribe(val => this.getActualTime());
    startEndTimes: StartEndTime = new StartEndTime(-1, "-", "-")
    getStartEndTime(): void {
        if (this.server_id === -1) {
            this.startEndTimes.start_time = "-";
            this.startEndTimes.end_time = "-";``
            return
        }
        let request = `${this.URLHead}/time/start&end?server_id=${this.server_id}`
        this.http.get(request).subscribe({next:(data:StartEndTime) => this.startEndTimes=data[0], error: error => console.error(error)});
    }
    ngOnDestroy() {
        this.subscription.unsubscribe();
    }    
}
