pragma solidity >= 0.7.0 < 0.9.0;

contract Play {

    mapping(string => uint256) public myDir;

    constructor(string memory _name, uint256 _mobNo) {
        myDir[_name] = _mobNo;
    }

    function setMobNo(string memory _name, uint256 _mobNo) public {
        myDir[_name] = _mobNo;
    }

    function getMobNo(string memory _name) public view returns (uint256) {
        return myDir[_name];
    }
}