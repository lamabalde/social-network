export class routerParams {
    constructor(location, param) {
        this.name = location
        this.params = {
            id : param
        }
    }
}