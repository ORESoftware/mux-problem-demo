


const v = [1,2,3];

const f = function(...values){
    console.log(values);   // [1,2,3,4]
};

f(...v, 4);