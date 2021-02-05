import sub2 from './sub1/sub1'
import root1 from '../../root-file-1'
import sub1 from '../function-level-file-1'
import sub3 from './sub1/sub2/sub2'
import sameName from '../same-name'

export function myFunc() {

}

function myFunc2() {

}


export default myFunc2

export const myFunc3 = () => {}

const myFunc4 = () => {}

export class MyClass1 {

}

export const myObjet = {}

export {
    myFunc2,
    myFunc4
}
